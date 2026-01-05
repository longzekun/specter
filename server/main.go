package main

import (
	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/command"
	"github.com/longzekun/specter/server/console"
	"github.com/longzekun/specter/server/constants"
	"github.com/longzekun/specter/server/db"
	"github.com/longzekun/specter/server/log"
	"go.uber.org/zap"
)

func main() {
	log.Init()
	certs.SetupCAs()

	db.Session()

	zap.S().Debugf("Starting Specter Server")
	con := console.NewConsole()

	server := con.Console.Menu(constants.ServerMenu)
	server.SetCommands(command.ServerCommands(con))

	con.Console.Start()
}
