package exit

import (
	"fmt"
	"os"

	"github.com/longzekun/specter/server/console"
	"github.com/longzekun/specter/server/constants"

	"github.com/spf13/cobra"
)

func ExitCmd(con *console.SpecterClient, cmd *cobra.Command, args []string) {
	fmt.Println("Exiting...")

	os.Exit(0)
}

func Command(con *console.SpecterClient) []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "exit",
			Short: "Exit the program",
			Run: func(cmd *cobra.Command, args []string) {
				ExitCmd(con, cmd, args)
			},
			GroupID: constants.GenericHelpGroup,
		},
	}
}
