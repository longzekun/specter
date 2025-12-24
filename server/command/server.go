package command

import (
	"github.com/longzekun/specter/server/command/exit"
	"github.com/longzekun/specter/server/command/jobs"
	client "github.com/longzekun/specter/server/console"
	"github.com/longzekun/specter/server/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func ServerCommands(con *client.SpecterClient) console.Commands {
	serverCommands := func() *cobra.Command {
		server := &cobra.Command{
			Short: "Server Commands",
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		}

		//	将服务端的命令挂载到server根上
		bind := makeBind(server, con)

		//	绑定通用命令
		bind(
			constants.GenericHelpGroup,
			exit.Command,
		)

		//	绑定网络命令
		bind(
			constants.NetworkHelpGroup,
			jobs.Command,
		)

		server.InitDefaultHelpCmd()
		return server
	}

	return serverCommands
}
