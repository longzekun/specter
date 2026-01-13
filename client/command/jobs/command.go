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

	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "control jobs",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ControlJobs(cmd, con, args)
		},
		GroupID: constants.NetworkHelpGroup,
	}
	flags.Bind("jobs control", false, jobsCmd, func(f *pflag.FlagSet) {
		f.Uint32P("kill", "k", 0, "kill job by id")
		f.BoolP("kill-all", "K", false, "kill all jobs")
	})

	return []*cobra.Command{
		mtlsCmd, jobsCmd,
	}
}
