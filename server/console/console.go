package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/longzekun/specter/server/constants"
	"github.com/reeflective/console"
	"github.com/reeflective/readline"
)

type SpecterClient struct {
	Console *console.Console
}

func NewConsole() *SpecterClient {
	c := &SpecterClient{
		Console: console.New("specter"),
	}

	//	通用配置
	c.Console.NewlineBefore = true
	c.Console.NewlineAfter = true

	//	add server control command menu
	server := c.Console.Menu(constants.ServerMenu)
	server.Short = "Server Command"
	server.Prompt().Primary = c.Prompt
	server.AddInterrupt(readline.ErrInterrupt, c.exitConsole)

	return c
}

func (c *SpecterClient) Prompt() string {
	//判断当前处于激活状态的是server还是implant,选择激活的menu
	return "specter >"
}

// exitConsole - 退出终端执行
func (c *SpecterClient) exitConsole(_ *console.Console) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y, Ctrl-C): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}
}
