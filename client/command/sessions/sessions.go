package sessions

import (
	"context"
	"fmt"
	"strings"

	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/spf13/cobra"
)

func ControlSessions(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	killId, _ := cmd.Flags().GetString("kill")
	if killId != "" {
		con.RPC.KillSession(context.Background(), &clientpb.KillReq{SessionID: killId})
		return
	}

	allSessions, err := con.RPC.GetAllSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		con.Printf("Failed to get all sessions: %v\n", err)
		return
	}

	con.Printf("%-40s %-10s %-10s %-24s %-10s %-10s %-70s %-24s %s\n",
		"ID",
		"Name",
		"Transport",
		"Remote Address",
		"Hostname",
		"Username",
		"Process (PID)",
		"Operating System",
		"Health",
	)
	con.Printf("%-40s %-10s %-10s %-24s %-10s %-10s %-70s %-24s %s\n",
		strings.Repeat("=", 40),
		strings.Repeat("=", 10),
		strings.Repeat("=", 10),
		strings.Repeat("=", 24),
		strings.Repeat("=", 10),
		strings.Repeat("=", 10),
		strings.Repeat("=", 70),
		strings.Repeat("=", 24),
		strings.Repeat("=", 8),
	)

	for _, session := range allSessions.Sessions {
		processInfo := fmt.Sprintf("%s(%s)", session.Filename, session.PID)
		osArch := fmt.Sprintf("%s/%s", session.OS, session.Arch)
		var health string
		if session.Health {
			health = "\x1b[32m[ALIVE]\x1b[0m"
		} else {
			health = "\x1b[31m[DEAD]\x1b[0m"
		}

		con.Printf("%-40s %-10s %-10s %-24s %-10s %-10s %-70s %-24s %s\n",
			session.ID,
			session.Name,
			session.TransportType,
			session.RemoteAddress,
			session.Hostname,
			session.Username,
			processInfo,
			osArch,
			health,
		)
	}
}
