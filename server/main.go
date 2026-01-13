package main

import (
	"github.com/longzekun/specter/server/cli"
	"github.com/longzekun/specter/server/log"
)

func main() {
	log.Init()
	cli.Execute()
}
