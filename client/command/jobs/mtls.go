package jobs

import (
	"context"

	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func StartMTLSListener(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	lhost, _ := cmd.Flags().GetString("lhost")
	lport, _ := cmd.Flags().GetUint32("lport")

	//	判断端口是否在使用中

	//	开启mtls监听
	zap.S().Debugf("Listening on %s:%d", lhost, lport)

	mtls, err := con.RPC.StartMTLSListener(context.Background(),
		&clientpb.MTLSListenerReq{
			Host: lhost,
			Port: lport,
		})
	if err != nil {
		zap.S().Warnf("Failed to start MTLS listener: %s", err.Error())
	} else {
		zap.S().Infof("MTLS job id is %v", mtls.JobID)
	}
}
