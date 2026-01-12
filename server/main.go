package main

import (
	"os"

	"github.com/longzekun/specter/server/cli"
	"github.com/longzekun/specter/server/log"
	"github.com/longzekun/specter/server/version"
)

func main() {
	log.Init()
	version.BuildSpecterInfo()
	os.Exit(1)
	cli.Execute()
}
