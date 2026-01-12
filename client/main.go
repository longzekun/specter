package main

import (
	"context"
	"fmt"
	"os"

	"github.com/longzekun/specter/client/command"
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/log"
	"github.com/longzekun/specter/client/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func main() {
	log.Init()
	zap.S().Infof("Starting Specter Client")

	//	rpc connect
	rpc, ln, err := transport.MtlsConnect()
	if err != nil {
		zap.S().Fatalf("MTLS connect err: %v", err)
		return
	}

	go handleConnectionLost(ln)

	//	create console
	con := console.NewConsole(false)

	//	set rpc,server commands,implant  commands
	console.StartClient(con, rpc, command.ServerCommands(con, nil), nil)
}

func handleConnectionLost(ln *grpc.ClientConn) {
	currentState := ln.GetState()
	if ln.WaitForStateChange(context.Background(), currentState) {
		newState := ln.GetState()
		if newState == connectivity.Idle {
			fmt.Println("\nLost connection to server. Exiting now.")
			os.Exit(1)
		}
	}
}
