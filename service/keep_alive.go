package service

import (
	"context"
	"github.com/akirasalvare/fdu-connect/log"
	"github.com/akirasalvare/fdu-connect/resolve"
	"time"
)

func KeepAlive(resolver *resolve.Resolver) {
	remoteUDPResolver, err := resolver.RemoteUDPResolver()
	if err != nil {
		log.Printf("KeepAlive: %s", err)
		panic("KeepAlive: " + err.Error())
	}

	for {
		_, err := remoteUDPResolver.LookupIP(context.Background(), "ip4", "www.baidu.com")
		if err != nil {
			log.Printf("KeepAlive: %s", err)
		} else {
			log.Printf("KeepAlive: OK")
		}

		time.Sleep(60 * time.Second)
	}
}
