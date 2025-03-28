package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/akirasalvare/fdu-connect/configs"
)

func getTOMLVal[T int | uint64 | string | bool](valPointer *T, defaultVal T) T {
	if valPointer == nil {
		return defaultVal
	} else {
		return *valPointer
	}
}

func parseTOMLConfig(configFile string, conf *configs.Config) error {
	var confTOML configs.ConfigTOML

	_, err := toml.DecodeFile(configFile, &confTOML)
	if err != nil {
		return errors.New("FDU Connect: error parsing the config file")
	}

	conf.ServerAddress = getTOMLVal(confTOML.ServerAddress, "stuvpn.fudan.edu.cn")
	conf.ServerPort = getTOMLVal(confTOML.ServerPort, 443)
	conf.Username = getTOMLVal(confTOML.Username, "")
	conf.Password = getTOMLVal(confTOML.Password, "")
	conf.TOTPSecret = getTOMLVal(confTOML.TOTPSecret, "")
	conf.CertFile = getTOMLVal(confTOML.CertFile, "")
	conf.CertPassword = getTOMLVal(confTOML.CertPassword, "")
	conf.DisableServerConfig = getTOMLVal(confTOML.DisableServerConfig, false)
	conf.SkipDomainResource = getTOMLVal(confTOML.SkipDomainResource, false)
	conf.DisableFDUConfig = getTOMLVal(confTOML.DisableFDUConfig, false)
	conf.DisableFDUDNS = getTOMLVal(confTOML.DisableFDUDNS, false)
	conf.DisableMultiLine = getTOMLVal(confTOML.DisableMultiLine, false)
	conf.ProxyAll = getTOMLVal(confTOML.ProxyAll, false)
	conf.SocksBind = getTOMLVal(confTOML.SocksBind, ":1080")
	conf.SocksUser = getTOMLVal(confTOML.SocksUser, "")
	conf.SocksPasswd = getTOMLVal(confTOML.SocksPasswd, "")
	conf.HTTPBind = getTOMLVal(confTOML.HTTPBind, ":1081")
	conf.ShadowsocksURL = getTOMLVal(confTOML.ShadowsocksURL, "")
	conf.DialDirectProxy = getTOMLVal(confTOML.DialDirectProxy, "")
	conf.TUNMode = getTOMLVal(confTOML.TUNMode, false)
	conf.AddRoute = getTOMLVal(confTOML.AddRoute, false)
	conf.DNSTTL = getTOMLVal(confTOML.DNSTTL, uint64(3600))
	conf.DebugDump = getTOMLVal(confTOML.DebugDump, false)
	conf.DisableKeepAlive = getTOMLVal(confTOML.DisableKeepAlive, false)
	conf.FDUDNSServer = getTOMLVal(confTOML.FDUDNSServer, "202.120.224.26")
	conf.SecondaryDNSServer = getTOMLVal(confTOML.SecondaryDNSServer, "223.5.5.5")
	conf.DNSServerBind = getTOMLVal(confTOML.DNSServerBind, "")
	conf.DNSHijack = getTOMLVal(confTOML.DNSHijack, false)

	for _, singlePortForwarding := range confTOML.PortForwarding {
		if singlePortForwarding.NetworkType == nil {
			return errors.New("FDU Connect: network type is not set")
		}

		if singlePortForwarding.BindAddress == nil {
			return errors.New("FDU Connect: bind address is not set")
		}

		if singlePortForwarding.RemoteAddress == nil {
			return errors.New("FDU Connect: remote address is not set")
		}

		conf.PortForwardingList = append(conf.PortForwardingList, configs.SinglePortForwarding{
			NetworkType:   *singlePortForwarding.NetworkType,
			BindAddress:   *singlePortForwarding.BindAddress,
			RemoteAddress: *singlePortForwarding.RemoteAddress,
		})
	}

	for _, singleCustomDns := range confTOML.CustomDNS {
		if singleCustomDns.HostName == nil {
			return errors.New("FDU Connect: host name is not set")
		}

		if singleCustomDns.IP == nil {
			fmt.Println("FDU Connect: IP is not set")
			return errors.New("FDU Connect: IP is not set")
		}

		conf.CustomDNSList = append(conf.CustomDNSList, configs.SingleCustomDNS{
			HostName: *singleCustomDns.HostName,
			IP:       *singleCustomDns.IP,
		})
	}

	for _, singleCustomProxyDomain := range confTOML.CustomProxyDomain {
		var domainRegex = regexp.MustCompile(`^[a-zA-Z\d-]+(\.[a-zA-Z\d-]+)*\.[a-zA-Z]{2,}$`)
		if !domainRegex.MatchString(singleCustomProxyDomain) {
			fmt.Printf("FDU Connect: %s is not a valid domain\n", singleCustomProxyDomain)
			return errors.New(fmt.Sprintf("FDU Connect: %s is not a valid domain", singleCustomProxyDomain))
		}
		conf.CustomProxyDomain = append(conf.CustomProxyDomain, singleCustomProxyDomain)
	}

	return nil
}

func init() {
	configFile, tcpPortForwarding, udpPortForwarding, customDns, customProxyDomain := "", "", "", "", ""
	showVersion := false

	flag.StringVar(&conf.ServerAddress, "server", "stuvpn.fudan.edu.cn", "EasyConnect server address")
	flag.IntVar(&conf.ServerPort, "port", 443, "EasyConnect port address")
	flag.StringVar(&conf.Username, "username", "", "Your username")
	flag.StringVar(&conf.Password, "password", "", "Your password")
	flag.StringVar(&conf.TOTPSecret, "totp-secret", "", "TOTP secret")
	flag.StringVar(&conf.CertFile, "cert-file", "", "Client certificate p12 file path for certificate login")
	flag.StringVar(&conf.CertPassword, "cert-password", "", "Client certificate password")
	flag.BoolVar(&conf.DisableServerConfig, "disable-server-config", false, "Don't parse server config")
	flag.BoolVar(&conf.SkipDomainResource, "skip-domain-resource", false, "Don't use server domain resource to decide whether to use RVPN.")
	flag.BoolVar(&conf.DisableFDUConfig, "disable-fdu-config", false, "Don't use FDU config")
	flag.BoolVar(&conf.DisableFDUDNS, "disable-fdu-dns", false, "Use local DNS instead of FDU DNS")
	flag.BoolVar(&conf.DisableMultiLine, "disable-multi-line", false, "Disable multi line auto select")
	flag.BoolVar(&conf.ProxyAll, "proxy-all", false, "Proxy all IPv4 traffic")
	flag.StringVar(&conf.SocksBind, "socks-bind", ":1080", "The address SOCKS5 server listens on (e.g. 127.0.0.1:1080)")
	flag.StringVar(&conf.SocksUser, "socks-user", "", "SOCKS5 username, default is don't use auth")
	flag.StringVar(&conf.SocksPasswd, "socks-passwd", "", "SOCKS5 password, default is don't use auth")
	flag.StringVar(&conf.HTTPBind, "http-bind", ":1081", "The address HTTP server listens on (e.g. 127.0.0.1:1081)")
	flag.StringVar(&conf.ShadowsocksURL, "shadowsocks-url", "", "The address Shadowsocks server listens on (e.g. ss://method:password@host:port)")
	flag.StringVar(&conf.DialDirectProxy, "dial-direct-proxy", "", "Dial with proxy when the connection doesn't match RVPN rules (e.g. http://127.0.0.1:7890)")
	flag.BoolVar(&conf.TUNMode, "tun-mode", false, "Enable TUN mode (experimental)")
	flag.BoolVar(&conf.AddRoute, "add-route", false, "Add route from rules for TUN interface")
	flag.Uint64Var(&conf.DNSTTL, "dns-ttl", 3600, "DNS record time to live, unit is second")
	flag.BoolVar(&conf.DebugDump, "debug-dump", false, "Enable traffic debug dump (only for debug usage)")
	flag.BoolVar(&conf.DisableKeepAlive, "disable-keep-alive", false, "Disable keep alive")
	flag.StringVar(&conf.FDUDNSServer, "fdu-dns-server", "202.120.224.26", "FDU DNS server address. Set to 'auto' to use DNS server provided by server")
	flag.StringVar(&conf.SecondaryDNSServer, "secondary-dns-server", "223.5.5.5", "Secondary DNS server address. Leave empty to use system default DNS server")
	flag.StringVar(&conf.DNSServerBind, "dns-server-bind", "", "The address DNS server listens on (e.g. 127.0.0.1:53)")
	flag.BoolVar(&conf.DNSHijack, "dns-hijack", false, "Hijack all dns query to FDU Connect")
	flag.StringVar(&conf.TwfID, "twf-id", "", "Login using twfID captured (mostly for debug usage)")
	flag.StringVar(&tcpPortForwarding, "tcp-port-forwarding", "", "TCP port forwarding (e.g. 0.0.0.0:9898-10.10.98.98:80,127.0.0.1:9899-10.10.98.98:80)")
	flag.StringVar(&udpPortForwarding, "udp-port-forwarding", "", "UDP port forwarding (e.g. 127.0.0.1:53-202.120.224.26:53)")
	flag.StringVar(&customDns, "custom-dns", "", "Custom set dns lookup (e.g. www.cc98.org:10.10.98.98,appservice.fdu.edu.cn:10.203.8.198)")
	flag.StringVar(&customProxyDomain, "custom-proxy-domain", "", "Custom set domains which force use RVPN proxy  (e.g. science.org, nature.com)")
	flag.StringVar(&configFile, "config", "", "Config file")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("FDU Connect v%s\n", fduConnectVersion)
		os.Exit(0)
	}

	if configFile != "" {
		err := parseTOMLConfig(configFile, &conf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if tcpPortForwarding != "" {
			forwardingStringList := strings.Split(tcpPortForwarding, ",")
			for _, forwardingString := range forwardingStringList {
				addressStringList := strings.Split(forwardingString, "-")
				if len(addressStringList) != 2 {
					fmt.Println("FDU Connect: wrong tcp port forwarding format")
					os.Exit(1)
				}

				conf.PortForwardingList = append(conf.PortForwardingList, configs.SinglePortForwarding{
					NetworkType:   "tcp",
					BindAddress:   addressStringList[0],
					RemoteAddress: addressStringList[1],
				})
			}
		}

		if udpPortForwarding != "" {
			forwardingStringList := strings.Split(udpPortForwarding, ",")
			for _, forwardingString := range forwardingStringList {
				addressStringList := strings.Split(forwardingString, "-")
				if len(addressStringList) != 2 {
					fmt.Println("FDU Connect: wrong udp port forwarding format")
					os.Exit(1)
				}

				conf.PortForwardingList = append(conf.PortForwardingList, configs.SinglePortForwarding{
					NetworkType:   "udp",
					BindAddress:   addressStringList[0],
					RemoteAddress: addressStringList[1],
				})
			}
		}

		if customDns != "" {
			dnsList := strings.Split(customDns, ",")
			for _, dnsString := range dnsList {
				dnsStringSplit := strings.Split(dnsString, ":")
				if len(dnsStringSplit) != 2 {
					fmt.Println("FDU Connect: wrong custom dns format")
					os.Exit(1)
				}

				conf.CustomDNSList = append(conf.CustomDNSList, configs.SingleCustomDNS{
					HostName: dnsStringSplit[0],
					IP:       dnsStringSplit[1],
				})
			}
		}

		if customProxyDomain != "" {
			domainList := strings.Split(customProxyDomain, ",")
			for _, domain := range domainList {
				var domainRegex = regexp.MustCompile(`^[a-zA-Z\d-]+(\.[a-zA-Z\d-]+)*\.[a-zA-Z]{2,}$`)
				if !domainRegex.MatchString(domain) {
					fmt.Printf("FDU Connect: %s is not a valid domain\n", domain)
					os.Exit(1)
				}
				conf.CustomProxyDomain = append(conf.CustomProxyDomain, domain)
			}
		}
	}

	if conf.ServerAddress == "" || ((conf.Username == "" || conf.Password == "") && conf.TwfID == "") {
		fmt.Println("FDU Connect")
		fmt.Println("Please see: https://github.com/akirasalvare/f di-connect")
		fmt.Printf("\nUsage: %s -username <username> -password <password>\n", os.Args[0])
		fmt.Println("\nFull usage:")
		flag.PrintDefaults()

		os.Exit(1)
	}
}
