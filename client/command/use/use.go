package use

import (
	"context"

	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/spf13/cobra"
)

func UseCmd(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	//var session *clientpb.Session
	var sessionId string
	if len(args) > 0 {
		sessionId = args[0]
	}

	if sessionId == "" {
		con.Printf("no session id specified")
		return
	}

	allSessions, err := con.RPC.GetAllSessions(context.Background(), &commonpb.Empty{})
	if err != nil {
		con.Printf("get all sessions failed: %v\n", err)
		return
	}

	var activeSession *clientpb.Session
	for _, session := range allSessions.Sessions {
		if session.ID == sessionId {
			activeSession = session
			break
		}
	}

	if activeSession != nil {
		con.Printf("Active session %s (%s)\n", activeSession.Name, activeSession.ID)
		con.ActiveTarget.SetSession(activeSession)
	}
}
