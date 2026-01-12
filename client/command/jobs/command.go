package jobs

import (
	"github.com/longzekun/specter/client/command/flags"
	"github.com/longzekun/specter/client/command/generate"
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command(con *console.SpecterClient) []*cobra.Command {
	mtlsCmd := &cobra.Command{
		Use:   constants.MtlsStr,
		Short: "Start an mTls listener",
		Run: func(cmd *cobra.Command, args []string) {
			StartMTLSListener(cmd, con, args)
		},
		GroupID: constants.NetworkHelpGroup,
	}
	flags.Bind("mTls listener", false, mtlsCmd, func(f *pflag.FlagSet) {
		f.StringP("lhost", "L", "", "interface to bind server to")
		f.Uint32P("lport", "l", generate.MTLSDefaultPort, "tcp listen port")
	})

	return []*cobra.Command{
		mtlsCmd,
	}
}
