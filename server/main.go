package main

import (
	"time"

	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/log"
	"github.com/longzekun/specter/server/transport"
	"go.uber.org/zap"
)

func main() {
	log.Init()

	certs.SetupCAs()

	zap.S().Debugf("Starting Specter Server")

	transport.StartMtlsClientListener("0.0.0.0", 7777)

	for {
		time.Sleep(1 * time.Second)
	}
}
