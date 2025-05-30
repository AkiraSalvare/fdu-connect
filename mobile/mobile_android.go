package mobile

import (
	"crypto/tls"
	"github.com/akirasalvare/fdu-connect/client"
	"github.com/akirasalvare/fdu-connect/log"
	"github.com/akirasalvare/fdu-connect/stack/tun"
)

var vpnClient *client.EasyConnectClient

func Login(server string, username string, password string) string {
	log.Init()

	vpnClient = client.NewEasyConnectClient(
		server,
		username,
		password,
		"",
		tls.Certificate{},
		"",
		false,
		false,
		false,
	)
	err := vpnClient.Setup()
	if err != nil {
		return ""
	}

	log.Printf("EasyConnect client started")

	clientIP, err := vpnClient.IP()
	if err != nil {
		return ""
	}

	return clientIP.String()
}

func DebugLogin(server string, username string, password string) string {
	log.Init()
	log.EnableDebug()

	vpnClient = client.NewEasyConnectClient(
		server,
		username,
		password,
		"",
		tls.Certificate{},
		"",
		false,
		false,
		false,
	)
	err := vpnClient.Setup()
	if err != nil {
		return ""
	}

	log.Printf("EasyConnect client started")

	clientIP, err := vpnClient.IP()
	if err != nil {
		return ""
	}

	return clientIP.String()
}

func StartStack(fd int) {
	vpnTUNStack, err := tun.NewStack(vpnClient)
	if err != nil {
		return
	}

	vpnTUNStack.SetupTun(fd)
	vpnTUNStack.Run()
}
