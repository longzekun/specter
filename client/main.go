package main

import (
	"github.com/longzekun/specter/client/command"
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/log"
	"github.com/longzekun/specter/client/transport"
	"go.uber.org/zap"
)

func main() {
	log.Init()
	zap.S().Infof("Starting Specter Client")

	//	rpc connect
	rpc, _, err := transport.MtlsConnect()
	if err != nil {
		zap.S().Fatalf("MTLS connect err: %v", err)
		return
	}

	//	create console
	con := console.NewConsole(false)

	//	set rpc,server commands,implant  commands
	console.StartClient(con, rpc, command.ServerCommands(con), nil)

}
