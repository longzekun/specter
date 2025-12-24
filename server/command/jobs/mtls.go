package jobs

import (
	"github.com/longzekun/specter/server/c2"
	"github.com/longzekun/specter/server/console"
	"github.com/spf13/cobra"
)

func StartMTLSListener(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	lhost, _ := cmd.Flags().GetString("lhost")
	lport, _ := cmd.Flags().GetUint32("lport")

	//	判断端口是否在使用中

	//	开启mtls监听
	c2.StartMTLSListenerJob(lhost, lport)
}
