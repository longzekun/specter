package main

import (
	"github.com/longzekun/specter/server/command"
	"github.com/longzekun/specter/server/console"
	"github.com/longzekun/specter/server/constants"
)

func main() {
	con := console.NewConsole()

	server := con.Console.Menu(constants.ServerMenu)
	server.SetCommands(command.ServerCommands(con))

	con.Console.Start()
}
