package sessions

import (
	"github.com/longzekun/specter/client/command/flags"
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Command(con *console.SpecterClient) []*cobra.Command {
	sessionsCmd := &cobra.Command{
		Use:   constants.SessionsStr,
		Short: "Session Management",
		Run: func(cmd *cobra.Command, args []string) {
			ControlSessions(cmd, con, args)
		},
		GroupID: constants.SpecterHelpGroup,
	}

	flags.Bind("Session Management", false, sessionsCmd, func(f *pflag.FlagSet) {
		f.StringP("kill", "k", "", "kill session by id")
	})

	return []*cobra.Command{
		sessionsCmd,
	}
}
