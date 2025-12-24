package command

import (
	client "github.com/longzekun/specter/server/console"
	"github.com/spf13/cobra"
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
