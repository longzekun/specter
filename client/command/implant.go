package command

import (
	"github.com/longzekun/specter/client/command/screenshot"
	client "github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func ImplantCommands(con *client.SpecterClient) console.Commands {
	implantCommands := func() *cobra.Command {
		implant := &cobra.Command{
			Short: "Implant commands",
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		}

		bind := makeBind(implant, con)

		//info
		bind(constants.ImplantInfoHelpGroup,
			screenshot.Command,
		)

		return implant
	}

	return implantCommands
}
