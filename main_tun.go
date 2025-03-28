//go:build tun

package main

import (
	"context"
	"crypto"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/akirasalvare/fdu-connect/client"
	"github.com/akirasalvare/fdu-connect/configs"
	"github.com/akirasalvare/fdu-connect/dial"
	"github.com/akirasalvare/fdu-connect/internal/hook_func"
	"github.com/akirasalvare/fdu-connect/log"
	"github.com/akirasalvare/fdu-connect/resolve"
	"github.com/akirasalvare/fdu-connect/service"
	"github.com/akirasalvare/fdu-connect/stack/tun"
	"golang.org/x/crypto/pkcs12"
	"inet.af/netaddr"
)

var conf configs.Config

const fduConnectVersion = "0.9.0-tun-only"

func main() {
	log.Init()

	log.Println("Start FDU Connect v" + fduConnectVersion)
	if conf.DebugDump {
		log.EnableDebug()
	}

	if errs := hook_func.ExecInitialFunc(context.Background(), conf); errs != nil {
		for _, err := range errs {
			log.Printf("Initial FDU-Connect failed: %s", err)
		}
		os.Exit(1)
	}

	tlsCert := tls.Certificate{}
	if conf.CertFile != "" {
		p12Data, err := os.ReadFile(conf.CertFile)
		if err != nil {
			log.Fatalf("Read certificate file error: %s", err)
		}

		key, cert, err := pkcs12.Decode(p12Data, conf.CertPassword)
		if err != nil {
			log.Fatalf("Decode certificate file error: %s", err)
		}

		tlsCert = tls.Certificate{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key.(crypto.PrivateKey),
			Leaf:        cert,
		}
	}

	vpnClient := client.NewEasyConnectClient(
		conf.ServerAddress+":"+fmt.Sprintf("%d", conf.ServerPort),
		conf.Username,
		conf.Password,
		conf.TOTPSecret,
		tlsCert,
		conf.TwfID,
		!conf.DisableMultiLine,
		!conf.DisableServerConfig,
		!conf.SkipDomainResource,
	)
	err := vpnClient.Setup()
	if err != nil {
		log.Fatalf("EasyConnect client setup error: %s", err)
	}

	log.Printf("EasyConnect client started")

	ipResources, err := vpnClient.IPResources()
	if err != nil && !conf.DisableServerConfig {
		log.Println("No IP resources")
	}

	ipSet, err := vpnClient.IPSet()
	if err != nil && !conf.DisableServerConfig {
		log.Println("No IP set")
	}

	domainResources, err := vpnClient.DomainResources()
	if err != nil && !conf.DisableServerConfig {
		log.Println("No domain resources")
	}

	dnsResource, err := vpnClient.DNSResource()
	if err != nil && !conf.DisableServerConfig {
		log.Println("No DNS resource")
	}

	if !conf.DisableFDUConfig {
		if domainResources != nil {
			domainResources["fudan.edu.cn"] = client.DomainResource{
				PortMin:  1,
				PortMax:  65535,
				Protocol: "all",
			}
		} else {
			domainResources = map[string]client.DomainResource{
				"fudan.edu.cn": {
					PortMin:  1,
					PortMax:  65535,
					Protocol: "all",
				},
			}
		}

		if ipResources != nil {
			ipResources = append(ipResources, client.IPResource{
				IPMin:    net.ParseIP("10.0.0.0"),
				IPMax:    net.ParseIP("10.255.255.255"),
				PortMin:  1,
				PortMax:  65535,
				Protocol: "all",
			})
		} else {
			ipResources = []client.IPResource{{
				IPMin:    net.ParseIP("10.0.0.0"),
				IPMax:    net.ParseIP("10.255.255.255"),
				PortMin:  1,
				PortMax:  65535,
				Protocol: "all",
			}}
		}

		ipSetBuilder := netaddr.IPSetBuilder{}
		if ipSet != nil {
			ipSetBuilder.AddSet(ipSet)
		}
		ipSetBuilder.AddPrefix(netaddr.MustParseIPPrefix("10.0.0.0/8"))
		ipSet, _ = ipSetBuilder.IPSet()
	}

	for _, customProxyDomain := range conf.CustomProxyDomain {
		if domainResources != nil {
			domainResources[customProxyDomain] = client.DomainResource{
				PortMin:  1,
				PortMax:  65535,
				Protocol: "all",
			}
		} else {
			domainResources = map[string]client.DomainResource{
				customProxyDomain: {
					PortMin:  1,
					PortMax:  65535,
					Protocol: "all",
				},
			}
		}
	}

	vpnStack, err := tun.NewStack(vpnClient, conf.DNSHijack, ipResources)
	if err != nil {
		log.Fatalf("Tun stack setup error, make sure you are root user : %s", err)
	}

	if conf.AddRoute && ipSet != nil {
		for _, prefix := range ipSet.Prefixes() {
			log.Printf("Add route to %s", prefix.String())
			_ = vpnStack.AddRoute(prefix.String())
		}
	} else if !conf.AddRoute && !conf.DisableFDUConfig {
		log.Println("Add route to 10.0.0.0/8")
		_ = vpnStack.AddRoute("10.0.0.0/8")
	}

	useFDUDNS := !conf.DisableFDUDNS
	fduDNSServer := conf.FDUDNSServer
	if useFDUDNS && fduDNSServer == "auto" {
		fduDNSServer, err = vpnClient.DNSServer()
		if err != nil {
			useFDUDNS = false
			fduDNSServer = "202.120.224.26"
			log.Println("No DNS server provided by server. Disable FDU DNS")
		} else {
			log.Printf("Use DNS server %s provided by server", fduDNSServer)
		}
	}

	vpnResolver := resolve.NewResolver(
		vpnStack,
		fduDNSServer,
		conf.SecondaryDNSServer,
		conf.DNSTTL,
		domainResources,
		dnsResource,
		useFDUDNS,
	)

	for _, customDns := range conf.CustomDNSList {
		ipAddr := net.ParseIP(customDns.IP)
		if ipAddr == nil {
			log.Printf("Custom DNS for host name %s is invalid, SKIP", customDns.HostName)
		}
		vpnResolver.SetPermanentDNS(customDns.HostName, ipAddr)
		log.Printf("Add custom DNS: %s -> %s\n", customDns.HostName, customDns.IP)
	}
	localResolver := service.NewDnsServer(vpnResolver, []string{fduDNSServer, conf.SecondaryDNSServer})
	vpnStack.SetupResolve(localResolver)

	go vpnStack.Run()

	vpnDialer := dial.NewDialer(vpnStack, vpnResolver, ipResources, conf.ProxyAll, conf.DialDirectProxy)

	if conf.DNSServerBind != "" {
		go service.ServeDNS(conf.DNSServerBind, localResolver)
	}
	clientIP, _ := vpnClient.IP()
	go service.ServeDNS(clientIP.String()+":53", localResolver)

	if conf.SocksBind != "" {
		go service.ServeSocks5(conf.SocksBind, vpnDialer, vpnResolver, conf.SocksUser, conf.SocksPasswd)
	}

	if conf.HTTPBind != "" {
		go service.ServeHTTP(conf.HTTPBind, vpnDialer)
	}

	if conf.ShadowsocksURL != "" {
		go service.ServeShadowsocks(vpnDialer, conf.ShadowsocksURL)
	}

	for _, portForwarding := range conf.PortForwardingList {
		if portForwarding.NetworkType == "tcp" {
			go service.ServeTCPForwarding(vpnStack, portForwarding.BindAddress, portForwarding.RemoteAddress)
		} else if portForwarding.NetworkType == "udp" {
			go service.ServeUDPForwarding(vpnStack, portForwarding.BindAddress, portForwarding.RemoteAddress)
		} else {
			log.Printf("Port forwarding: unknown network type %s. Aborting", portForwarding.NetworkType)
		}
	}

	if !conf.DisableKeepAlive {
		if !useFDUDNS {
			log.Println("Keep alive is disabled because FDU DNS is disabled")
		} else {
			go service.KeepAlive(vpnResolver)
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	log.Println("Shutdown FDU-Connect ......")
	if errs := hook_func.ExecTerminalFunc(context.Background()); errs != nil {
		for _, err := range errs {
			log.Printf("Shutdown FDU-Connect failed: %s", err)
		}
	} else {
		log.Println("Shutdown FDU-Connect success, Bye~")
	}
}
