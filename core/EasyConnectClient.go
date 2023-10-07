package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/mythologyli/zju-connect/core/config"
	"github.com/mythologyli/zju-connect/parser"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type Forwarding struct {
	NetworkType   string
	BindAddress   string
	RemoteAddress string
}

type CustomDNS struct {
	HostName string
	IP       string
}

var SocksBind string
var SocksUser string
var SocksPasswd string
var HttpBind string
var DebugDump bool
var ParseServConfig bool
var ParseZjuConfig bool
var UseZjuDns bool
var TestMultiLine bool
var DnsTTL uint64
var ProxyAll bool
var ForwardingList []Forwarding
var EnableKeepAlive bool
var ZjuDnsServer string
var CustomDNSList []CustomDNS

type EasyConnectClient struct {
	queryConn net.Conn
	clientIp  []byte
	token     *[48]byte
	twfId     string

	endpoint *EasyConnectEndpoint
	ipStack  *stack.Stack

	server   string
	username string
	password string
}

func NewEasyConnectClient(server string) *EasyConnectClient {
	return &EasyConnectClient{
		server: server,
	}
}

func StartClient(host string, port int, username string, password string, twfId string) {
	server := fmt.Sprintf("%s:%d", host, port)

	client := NewEasyConnectClient(server)

	var ip []byte
	var err error
	if twfId != "" {
		if len(twfId) != 16 {
			panic("len(twfid) should be 16!")
		}
		err = client.LoginByTwfId(twfId)
	} else {
		err = client.Login(username, password)
		if err == ERR_NEXT_AUTH_SMS {
			fmt.Print(">>>Please enter your sms code<<<:")
			smsCode := ""
			_, _err := fmt.Scan(&smsCode)
			if _err != nil {
				panic(_err)
			}

			err = client.AuthSMSCode(smsCode)
		} else if err == ERR_NEXT_AUTH_TOTP {
			fmt.Print(">>>Please enter your TOTP Auth code<<<:")
			TOTPCode := ""
			_, _err := fmt.Scan(&TOTPCode)
			if _err != nil {
				panic(_err)
			}

			err = client.AuthTOTP(TOTPCode)
		}
	}

	if TestMultiLine && config.IsServerListAvailable() {
		log.Printf("Testing %d servers...", config.GetServerListLen())

		server := config.GetBestServer()

		if server != "" {
			log.Printf("Find best server: %s", server)

			TestMultiLine = false

			parts := strings.Split(server, ":")
			host := parts[0]
			port, _ := strconv.Atoi(parts[1])

			log.Printf("Login again...")

			StartClient(host, port, username, password, twfId)
			return
		} else {
			log.Printf("Find best server failed. Connect %s", client.server)
		}
	}
	// check error after trying best server
	if err != nil {
		log.Fatal(err.Error())
	}

	client.ParseAllConfig()
	if ip, err = client.GetClientIp(); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Login success, your IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	}

	// Link-level endpoint used in gvisor netstack
	client.endpoint = &EasyConnectEndpoint{}
	client.ipStack = SetupStack(client.clientIp, client.endpoint)

	// Sangfor Easyconnect protocol
	StartProtocol(client.endpoint, client.server, client.token, &[4]byte{client.clientIp[3], client.clientIp[2], client.clientIp[1], client.clientIp[0]}, DebugDump)

	for _, singleForwarding := range ForwardingList {
		go client.ServeForwarding(strings.ToLower(singleForwarding.NetworkType), singleForwarding.BindAddress, singleForwarding.RemoteAddress)
	}

	for _, customDNS := range CustomDNSList {
		ipAddr := net.ParseIP(customDNS.IP)
		if ipAddr == nil {
			log.Printf("Custom DNS for host_name %s is invalid, SKIP", customDNS.HostName)
		}
		SetPermantDns(customDNS.HostName, ipAddr)
		log.Printf("Custom DNS %s -> %s\n", customDNS.HostName, customDNS.IP)
	}

	if SocksBind != "" {
		go client.ServeSocks5(SocksBind, ZjuDnsServer)

		if HttpBind != "" {
			go client.ServeHttp(HttpBind, SocksBind, SocksUser, SocksPasswd)
		}
	}

	if EnableKeepAlive {
		go client.KeepAlive(ZjuDnsServer)
	}

	select {}
}

func (client *EasyConnectClient) Login(username string, password string) error {
	client.username = username
	client.password = password

	// Web login part (Get TWFID & ECAgent Token => Final token used in binary stream)
	twfId, err := WebLogin(client.server, client.username, client.password)

	// Store TWFID for AuthSMS
	client.twfId = twfId
	if err != nil {
		return err
	}

	return client.LoginByTwfId(twfId)
}

func (client *EasyConnectClient) AuthSMSCode(code string) error {
	if client.twfId == "" {
		return errors.New("SMS Auth not required")
	}

	twfId, err := AuthSms(client.server, client.username, client.password, client.twfId, code)
	if err != nil {
		return err
	}

	return client.LoginByTwfId(twfId)
}

func (client *EasyConnectClient) AuthTOTP(code string) error {
	if client.twfId == "" {
		return errors.New("TOTP Auth not required")
	}

	twfId, err := TOTPAuth(client.server, client.username, client.password, client.twfId, code)
	if err != nil {
		return err
	}

	return client.LoginByTwfId(twfId)
}

func (client *EasyConnectClient) LoginByTwfId(twfId string) error {
	agentToken, err := ECAgentToken(client.server, twfId)
	if err != nil {
		return err
	}

	parser.ParseConfLists(client.server, twfId, DebugDump)

	client.token = (*[48]byte)([]byte(agentToken + twfId))
	return nil
}

func (client *EasyConnectClient) ParseAllConfig() {
	// Parse Server config
	if ParseServConfig {
		parser.ParseResourceLists(client.server, client.twfId, DebugDump)
	}

	// Parse ZJU config
	if ParseZjuConfig {
		parser.ParseZjuDnsRules(DebugDump)
		parser.ParseZjuIpv4Rules(DebugDump)
		parser.ParseZjuForceProxyRules(DebugDump)
	}
}

func (client *EasyConnectClient) GetClientIp() ([]byte, error) {
	var err error
	// Query IP (keep the connection used, so it's not closed too early, otherwise i/o stream will be closed)
	client.clientIp, client.queryConn, err = QueryIp(client.server, client.token, DebugDump)
	if err != nil {
		return nil, err
	}

	return client.clientIp, nil
}

func (client *EasyConnectClient) ServeSocks5(socksBind string, dnsServer string) {
	ServeSocks5(client.ipStack, client.clientIp, socksBind, dnsServer)
}

func (client *EasyConnectClient) ServeHttp(httpBind string, socksBind string, socksUser string, socksPasswd string) {
	ServeHttp(httpBind, socksBind, socksUser, socksPasswd)
}

func (client *EasyConnectClient) ServeForwarding(networkType string, bindAddress string, remoteAddress string) {
	if networkType == "tcp" {
		log.Printf("Port forwarding (tcp): %s <- %s", bindAddress, remoteAddress)

		ServeTcpForwarding(bindAddress, remoteAddress, client.ipStack, client.clientIp)
	} else if networkType == "udp" {
		log.Printf("Port forwarding (udp): %s <- %s", bindAddress, remoteAddress)

		ServeUdpForwarding(bindAddress, remoteAddress, client.ipStack)
	} else {
		log.Println("Only TCP/UDP forwarding is supported yet. Aborting.")
	}
}

func (client *EasyConnectClient) KeepAlive(dnsServer string) {
	KeepAlive(dnsServer, client.ipStack, client.clientIp)
}
