package console

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/longzekun/specter/client/constants"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/reeflective/console"
	"github.com/reeflective/readline"
)

type SpecterClient struct {
	Console  *console.Console
	IsServer bool
	RPC      rpcpb.SpecterRPCClient
	printf   func(format string, args ...any) (int, error)
}

func NewConsole(isServer bool) *SpecterClient {
	c := &SpecterClient{
		Console:  console.New("specter"),
		IsServer: isServer,
	}

	//	通用配置
	c.Console.NewlineBefore = true
	c.Console.NewlineAfter = true

	//	add server control command menu
	server := c.Console.Menu(constants.ServerMenu)
	server.Short = "Server Command"
	server.Prompt().Primary = c.Prompt
	server.AddInterrupt(readline.ErrInterrupt, c.exitConsole)

	c.Console.SetPrintLogo(func(_ *console.Console) {
		c.PrintLogo()
	})

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

func StartClient(con *SpecterClient, rpc rpcpb.SpecterRPCClient, serverCommands console.Commands, implantCommands console.Commands) error {
	con.RPC = rpc

	con.printf = con.Console.Printf

	go con.StartEventLoop()

	if serverCommands != nil {
		server := con.Console.Menu(constants.ServerMenu)
		server.SetCommands(serverCommands)
	}

	if implantCommands != nil {
		implant := con.Console.Menu(constants.ImplantMenu)
		implant.SetCommands(implantCommands)
	}

	return con.Console.Start()
}

func (c *SpecterClient) StartEventLoop() {
	eventStream, err := c.RPC.Events(context.Background(), &commonpb.Empty{})
	if err != nil {
		return
	}

	for {
		event, err := eventStream.Recv()
		if err == io.EOF || event == nil {
			return
		}

		//	deal event
		switch event.EventType {
		case constants.ClientJoinType:
			c.printf("\nClient joined,Operator name is %v\n", event.Client.Operator.Name)
		case constants.ClientLeaveType:
			c.printf("\nClient left,Operator name is %v\n", event.Client.Operator.Name)
		}

	}
}

func (c *SpecterClient) PrintLogo() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(7)
	fmt.Printf("%s", asciiSpecterLogo[n])

	//
}

const (
	ColorReset  = "\033[0m"  // 重置
	ColorRed    = "\033[31m" // 红色
	ColorGreen  = "\033[32m" // 绿色
	ColorYellow = "\033[33m" // 黄色
	ColorBlue   = "\033[34m" // 蓝色
	ColorPurple = "\033[35m" // 紫色
	ColorCyan   = "\033[36m" // 青色
	ColorWhite  = "\033[37m" // 白色
)

var asciiSpecterLogo = []string{
	ColorPurple +
		`
 ░▒▓███████▓▒░▒▓███████▓▒░░▒▓████████▓▒░▒▓██████▓▒░▒▓████████▓▒░▒▓████████▓▒░▒▓███████▓▒░
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░     ░▒▓█▓▒░        ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
 ░▒▓██████▓▒░░▒▓███████▓▒░░▒▓██████▓▒░░▒▓█▓▒░        ░▒▓█▓▒░   ░▒▓██████▓▒░ ░▒▓███████▓▒░
       ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░        ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
       ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓███████▓▒░░▒▓█▓▒░      ░▒▓████████▓▒░▒▓██████▓▒░  ░▒▓█▓▒░   ░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░
` + ColorReset,
	ColorRed +
		`
  ____     ____   U _____ u   ____   _____  U _____ u   ____     
 / __"| uU|  _"\ u\| ___"|/U /"___| |_ " _| \| ___"|/U |  _"\ u  
<\___ \/ \| |_) |/ |  _|"  \| | u     | |    |  _|"   \| |_) |/  
 u___) |  |  __/   | |___   | |/__   /| |\   | |___    |  _ <    
 |____/>> |_|      |_____|   \____| u |_|U   |_____|   |_| \_\   
  )(  (__)||>>_    <<   >>  _// \\  _// \\_  <<   >>   //   \\_  
 (__)    (__)__)  (__) (__)(__)(__)(__) (__)(__) (__) (__)  (__) 
` + ColorReset,
	ColorGreen + `
     _____        _____        _____        _____         _____        _____        _____     
  __|___  |__  __|__   |__  __|___  |__  __|___  |__  ___|__   |__  __|___  |__  __|__   |__  
 |   ___|    ||     |     ||   ___|    ||   ___|    ||_    _|     ||   ___|    ||     |     | 
  '-.'-.     ||    _|     ||   ___|    ||   |__     | |    |      ||   ___|    ||     \     |
 |______|  __||___|     __||______|  __||______|  __| |____|    __||______|  __||__|\__\  __|
    |_____|      |_____|      |_____|      |_____|       |_____|      |_____|      |_____|

` + ColorReset,
	ColorYellow + `
.------..------..------..------..------..------..------.
|S.--. ||P.--. ||E.--. ||C.--. ||T.--. ||E.--. ||R.--. |
| :/\: || :/\: || (\/) || :/\: || :/\: || (\/) || :(): |
| :\/: || (__) || :\/: || :\/: || (__) || :\/: || ()() |
| '--'S|| '--'P|| '--'E|| '--'C|| '--'T|| '--'E|| '--'R|
~------'~------'~------'~------'~------'~------'~------'
` + ColorReset,
	ColorBlue + `
   ___      ___    ___     ___    _____    ___     ___   
  / __|    | _ \  | __|   / __|  |_   _|  | __|   | _ \  
  \__ \    |  _/  | _|   | (__     | |    | _|    |   /  
  |___/   _|_|_   |___|   \___|   _|_|_   |___|   |_|_\  
_|"""""|_| """ |_|"""""|_|"""""|_|"""""|_|"""""|_|"""""| 
"'-0-0-'"'-0-0-'"'-0-0-'"'-0-0-'"'-0-0-'"'-0-0-'"'-0-0-' 
` + ColorReset,
	ColorCyan + ` 
      _/_/_/  _/_/_/    _/_/_/_/    _/_/_/  _/_/_/_/_/  _/_/_/_/  _/_/_/    
   _/        _/    _/  _/        _/            _/      _/        _/    _/   
    _/_/    _/_/_/    _/_/_/    _/            _/      _/_/_/    _/_/_/      
       _/  _/        _/        _/            _/      _/        _/    _/     
_/_/_/    _/        _/_/_/_/    _/_/_/      _/      _/_/_/_/  _/    _/       
` + ColorReset,
	` 
   _|_|_|  _|_|_|    _|_|_|_|    _|_|_|  _|_|_|_|_|  _|_|_|_|  _|_|_|    
 _|        _|    _|  _|        _|            _|      _|        _|    _|  
   _|_|    _|_|_|    _|_|_|    _|            _|      _|_|_|    _|_|_|    
       _|  _|        _|        _|            _|      _|        _|    _|  
 _|_|_|    _|        _|_|_|_|    _|_|_|      _|      _|_|_|_|  _|    _|   
`,
}
