package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

type (
	Config struct {
		ServerAddress       string
		ServerPort          int
		Username            string
		Password            string
		DisableServerConfig bool
		DisableZJUConfig    bool
		DisableZJUDNS       bool
		DisableMultiLine    bool
		ProxyAll            bool
		SocksBind           string
		SocksUser           string
		SocksPasswd         string
		HTTPBind            string
		TUNMode             bool
		AddRoute            bool
		DNSTTL              uint64
		DisableKeepAlive    bool
		ZJUDNSServer        string
		SecondaryDNSServer  string
		DNSServerBind       string
		TUNDNSServer        string
		DebugDump           bool
		PortForwardingList  []SinglePortForwarding
		CustomDNSList       []SingleCustomDNS
		TwfID               string
	}

	SinglePortForwarding struct {
		NetworkType   string
		BindAddress   string
		RemoteAddress string
	}

	SingleCustomDNS struct {
		HostName string `toml:"host_name"`
		IP       string `toml:"ip"`
	}
)

type (
	ConfigTOML struct {
		ServerAddress       *string                    `toml:"server_address"`
		ServerPort          *int                       `toml:"server_port"`
		Username            *string                    `toml:"username"`
		Password            *string                    `toml:"password"`
		DisableServerConfig *bool                      `toml:"disable_server_config"`
		DisableZJUConfig    *bool                      `toml:"disable_zju_config"`
		DisableZJUDNS       *bool                      `toml:"disable_zju_dns"`
		DisableMultiLine    *bool                      `toml:"disable_multi_line"`
		ProxyAll            *bool                      `toml:"proxy_all"`
		SocksBind           *string                    `toml:"socks_bind"`
		SocksUser           *string                    `toml:"socks_user"`
		SocksPasswd         *string                    `toml:"socks_passwd"`
		HTTPBind            *string                    `toml:"http_bind"`
		TUNMode             *bool                      `toml:"tun_mode"`
		AddRoute            *bool                      `toml:"add_route"`
		DNSTTL              *uint64                    `toml:"dns_ttl"`
		DisableKeepAlive    *bool                      `toml:"disable_keep_alive"`
		ZJUDNSServer        *string                    `toml:"zju_dns_server"`
		SecondaryDNSServer  *string                    `toml:"secondary_dns_server"`
		DNSServerBind       *string                    `toml:"dns_server_bind"`
		TUNDNSServer        *string                    `toml:"tun_dns_server"`
		DebugDump           *bool                      `toml:"debug_dump"`
		PortForwarding      []SinglePortForwardingTOML `toml:"port_forwarding"`
		CustomDNS           []SingleCustomDNSTOML      `toml:"custom_dns"`
	}

	SinglePortForwardingTOML struct {
		NetworkType   *string `toml:"network_type"`
		BindAddress   *string `toml:"bind_address"`
		RemoteAddress *string `toml:"remote_address"`
	}

	SingleCustomDNSTOML struct {
		HostName *string `toml:"host_name"`
		IP       *string `toml:"ip"`
	}
)

func getTOMLVal[T int | uint64 | string | bool](valPointer *T, defaultVal T) T {
	if valPointer == nil {
		return defaultVal
	} else {
		return *valPointer
	}
}

func parseTOMLConfig(configFile string, conf *Config) error {
	var confTOML ConfigTOML

	_, err := toml.DecodeFile(configFile, &confTOML)
	if err != nil {
		return errors.New("ZJU Connect: error parsing the config file")
	}

	conf.ServerAddress = getTOMLVal(confTOML.ServerAddress, "rvpn.zju.edu.cn")
	conf.ServerPort = getTOMLVal(confTOML.ServerPort, 443)
	conf.Username = getTOMLVal(confTOML.Username, "")
	conf.Password = getTOMLVal(confTOML.Password, "")
	conf.DisableServerConfig = getTOMLVal(confTOML.DisableServerConfig, false)
	conf.DisableZJUConfig = getTOMLVal(confTOML.DisableZJUConfig, false)
	conf.DisableZJUDNS = getTOMLVal(confTOML.DisableZJUDNS, false)
	conf.DisableMultiLine = getTOMLVal(confTOML.DisableMultiLine, false)
	conf.ProxyAll = getTOMLVal(confTOML.ProxyAll, false)
	conf.SocksBind = getTOMLVal(confTOML.SocksBind, ":1080")
	conf.SocksUser = getTOMLVal(confTOML.SocksUser, "")
	conf.SocksPasswd = getTOMLVal(confTOML.SocksPasswd, "")
	conf.HTTPBind = getTOMLVal(confTOML.HTTPBind, ":1081")
	conf.TUNMode = getTOMLVal(confTOML.TUNMode, false)
	conf.AddRoute = getTOMLVal(confTOML.AddRoute, false)
	conf.DNSTTL = getTOMLVal(confTOML.DNSTTL, uint64(3600))
	conf.DebugDump = getTOMLVal(confTOML.DebugDump, false)
	conf.DisableKeepAlive = getTOMLVal(confTOML.DisableKeepAlive, false)
	conf.ZJUDNSServer = getTOMLVal(confTOML.ZJUDNSServer, "10.10.0.21")
	conf.SecondaryDNSServer = getTOMLVal(confTOML.SecondaryDNSServer, "114.114.114.114")
	conf.DNSServerBind = getTOMLVal(confTOML.DNSServerBind, "")
	conf.TUNDNSServer = getTOMLVal(confTOML.TUNDNSServer, "")

	for _, singlePortForwarding := range confTOML.PortForwarding {
		if singlePortForwarding.NetworkType == nil {
			return errors.New("ZJU Connect: network type is not set")
		}

		if singlePortForwarding.BindAddress == nil {
			return errors.New("ZJU Connect: bind address is not set")
		}

		if singlePortForwarding.RemoteAddress == nil {
			return errors.New("ZJU Connect: remote address is not set")
		}

		conf.PortForwardingList = append(conf.PortForwardingList, SinglePortForwarding{
			NetworkType:   *singlePortForwarding.NetworkType,
			BindAddress:   *singlePortForwarding.BindAddress,
			RemoteAddress: *singlePortForwarding.RemoteAddress,
		})
	}

	for _, singleCustomDns := range confTOML.CustomDNS {
		if singleCustomDns.HostName == nil {
			return errors.New("ZJU Connect: host name is not set")
		}

		if singleCustomDns.IP == nil {
			fmt.Println("ZJU Connect: IP is not set")
			return errors.New("ZJU Connect: IP is not set")
		}

		conf.CustomDNSList = append(conf.CustomDNSList, SingleCustomDNS{
			HostName: *singleCustomDns.HostName,
			IP:       *singleCustomDns.IP,
		})
	}

	return nil
}

func init() {
	configFile, tcpPortForwarding, udpPortForwarding, customDns := "", "", "", ""
	showVersion := false

	flag.StringVar(&conf.ServerAddress, "server", "rvpn.zju.edu.cn", "EasyConnect server address")
	flag.IntVar(&conf.ServerPort, "port", 443, "EasyConnect port address")
	flag.StringVar(&conf.Username, "username", "", "Your username")
	flag.StringVar(&conf.Password, "password", "", "Your password")
	flag.BoolVar(&conf.DisableServerConfig, "disable-server-config", false, "Don't parse server config")
	flag.BoolVar(&conf.DisableZJUConfig, "disable-zju-config", false, "Don't use ZJU config")
	flag.BoolVar(&conf.DisableZJUDNS, "disable-zju-dns", false, "Use local DNS instead of ZJU DNS")
	flag.BoolVar(&conf.DisableMultiLine, "disable-multi-line", false, "Disable multi line auto select")
	flag.BoolVar(&conf.ProxyAll, "proxy-all", false, "Proxy all IPv4 traffic")
	flag.StringVar(&conf.SocksBind, "socks-bind", ":1080", "The address SOCKS5 server listens on (e.g. 127.0.0.1:1080)")
	flag.StringVar(&conf.SocksUser, "socks-user", "", "SOCKS5 username, default is don't use auth")
	flag.StringVar(&conf.SocksPasswd, "socks-passwd", "", "SOCKS5 password, default is don't use auth")
	flag.StringVar(&conf.HTTPBind, "http-bind", ":1081", "The address HTTP server listens on (e.g. 127.0.0.1:1081)")
	flag.BoolVar(&conf.TUNMode, "tun-mode", false, "Enable TUN mode (experimental)")
	flag.BoolVar(&conf.AddRoute, "add-route", false, "Add route from rules for TUN interface")
	flag.Uint64Var(&conf.DNSTTL, "dns-ttl", 3600, "DNS record time to live, unit is second")
	flag.BoolVar(&conf.DebugDump, "debug-dump", false, "Enable traffic debug dump (only for debug usage)")
	flag.BoolVar(&conf.DisableKeepAlive, "disable-keep-alive", false, "Disable keep alive")
	flag.StringVar(&conf.ZJUDNSServer, "zju-dns-server", "10.10.0.21", "ZJU DNS server address")
	flag.StringVar(&conf.SecondaryDNSServer, "secondary-dns-server", "114.114.114.114", "Secondary DNS server address. Leave empty to use system default DNS server")
	flag.StringVar(&conf.DNSServerBind, "dns-server-bind", "", "The address DNS server listens on (e.g. 127.0.0.1:53)")
	flag.StringVar(&conf.TUNDNSServer, "tun-dns-server", "", "DNS Server address for TUN interface (e.g. 127.0.0.1). You should not specify the port")
	flag.StringVar(&conf.TwfID, "twf-id", "", "Login using twfID captured (mostly for debug usage)")
	flag.StringVar(&tcpPortForwarding, "tcp-port-forwarding", "", "TCP port forwarding (e.g. 0.0.0.0:9898-10.10.98.98:80,127.0.0.1:9899-10.10.98.98:80)")
	flag.StringVar(&udpPortForwarding, "udp-port-forwarding", "", "UDP port forwarding (e.g. 127.0.0.1:53-10.10.0.21:53)")
	flag.StringVar(&customDns, "custom-dns", "", "Custom set dns lookup (e.g. www.cc98.org:10.10.98.98,appservice.zju.edu.cn:10.203.8.198)")
	flag.StringVar(&configFile, "config", "", "Config file")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("ZJU Connect v%s\n", zjuConnectVersion)
		return
	}

	if configFile != "" {
		err := parseTOMLConfig(configFile, &conf)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if tcpPortForwarding != "" {
			forwardingStringList := strings.Split(tcpPortForwarding, ",")
			for _, forwardingString := range forwardingStringList {
				addressStringList := strings.Split(forwardingString, "-")
				if len(addressStringList) != 2 {
					fmt.Println("ZJU Connect: wrong tcp port forwarding format")
					return
				}

				conf.PortForwardingList = append(conf.PortForwardingList, SinglePortForwarding{
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
					fmt.Println("ZJU Connect: wrong udp port forwarding format")
					return
				}

				conf.PortForwardingList = append(conf.PortForwardingList, SinglePortForwarding{
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
					fmt.Println("ZJU Connect: wrong custom dns format")
					return
				}

				conf.CustomDNSList = append(conf.CustomDNSList, SingleCustomDNS{
					HostName: dnsStringSplit[0],
					IP:       dnsStringSplit[1],
				})
			}
		}
	}

	if conf.ServerAddress == "" || ((conf.Username == "" || conf.Password == "") && conf.TwfID == "") {
		fmt.Println("ZJU Connect")
		fmt.Println("Please see: https://github.com/mythologyli/zju-connect")
		fmt.Printf("\nUsage: %s -username <username> -password <password>\n", os.Args[0])
		fmt.Println("\nFull usage:")
		flag.PrintDefaults()

		return
	}
}
