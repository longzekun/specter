package command

import (
	"github.com/longzekun/specter/client/command/exit"
	"github.com/longzekun/specter/client/command/jobs"
	"github.com/longzekun/specter/client/command/sessions"
	"github.com/longzekun/specter/client/command/use"
	client "github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func ServerCommands(con *client.SpecterClient, serverOnlyCommands func() []*cobra.Command) console.Commands {
	serverCommands := func() *cobra.Command {
		server := &cobra.Command{
			Short: "Server Commands",
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		}

		//	将服务端的命令挂载到server根上
		bind := makeBind(server, con)

		if serverOnlyCommands != nil {
			server.AddGroup(&cobra.Group{ID: constants.MultiplayerHelpGroup, Title: constants.MultiplayerHelpGroup})
			server.AddCommand(serverOnlyCommands()...)
		}

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

		bind(constants.SpecterHelpGroup,
			sessions.Command,
			use.Command,
		)

		server.InitDefaultHelpCmd()
		return server
	}

	return serverCommands
}
