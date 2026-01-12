package command

import (
	client "github.com/longzekun/specter/client/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func makeBind(cmd *cobra.Command, con *client.SpecterClient) func(group string, cmds ...func(con *client.SpecterClient) []*cobra.Command) {
	return func(group string, cmds ...func(con *client.SpecterClient) []*cobra.Command) {
		found := false

		if group != "" {
			for _, grp := range cmd.Groups() {
				if grp.Title == group {
					found = true
					break
				}
			}

			if !found {
				cmd.AddGroup(&cobra.Group{
					ID:    group,
					Title: group,
				})
			}
		}

		for _, command := range cmds {
			cmd.AddCommand(command(con)...)
		}
	}
}

func Bind(name string, persistent bool, cmd *cobra.Command, flags func(f *pflag.FlagSet)) {
	flagSet := pflag.NewFlagSet(name, pflag.ContinueOnError) // Create the flag set.
	flags(flagSet)                                           // Let the user bind any number of flags to it.

	if persistent {
		cmd.PersistentFlags().AddFlagSet(flagSet)
	} else {
		cmd.Flags().AddFlagSet(flagSet)
	}
}
