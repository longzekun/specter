package screenshot

import (
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/spf13/cobra"
)

func Command(con *console.SpecterClient) []*cobra.Command {
	screenshotCmd := &cobra.Command{
		Use:   constants.ScreenShotStr,
		Short: "Get implant screenshot",
		Run: func(cmd *cobra.Command, args []string) {
			ScreenshotCmd(cmd, con, args)
		},
		GroupID: constants.ImplantInfoHelpGroup,
	}

	return []*cobra.Command{
		screenshotCmd,
	}
}
