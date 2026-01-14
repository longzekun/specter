package screenshot

import (
	"github.com/longzekun/specter/client/console"
	"github.com/spf13/cobra"
)

func ScreenshotCmd(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	session := con.ActiveTarget.GetSession()
	if session == nil {
		con.Printf("no active session\n")
		return
	}
	con.Printf("save screenshot at ")
}
