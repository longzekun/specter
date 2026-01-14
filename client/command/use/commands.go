package use

import (
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/spf13/cobra"
)

func Command(con *console.SpecterClient) []*cobra.Command {
	useCmd := &cobra.Command{
		Use:   constants.UseStr,
		Short: "Switch the activate session",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			UseCmd(cmd, con, args)
		},
		GroupID: constants.SpecterHelpGroup,
	}

	return []*cobra.Command{useCmd}
}
