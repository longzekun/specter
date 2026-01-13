package jobs

import (
	"context"

	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/spf13/cobra"
)

func ControlJobs(cmd *cobra.Command, con *console.SpecterClient, args []string) {
	killId, _ := cmd.Flags().GetUint32("kill")
	killAll, _ := cmd.Flags().GetBool("kill-all")

	if killAll {
		con.RPC.KillAllJobs(context.Background(), &commonpb.Empty{})
		return
	}

	if killId != 0 {
		con.RPC.KillJob(context.Background(), &clientpb.KillJobReq{ID: killId})
		return
	}

	//	get all jobs
	jobs, err := con.RPC.GetAllJobs(context.Background(), &commonpb.Empty{})
	if err != nil {
		con.Printf("Failed to get all jobs: %v\n", err)
	}

	con.Printf("%-4s %-16s %-10s %-6s %s\n",
		"ID", "Name", "Protocol", "Port", "Description")
	con.Printf("%-4s %-16s %-10s %-6s %s\n",
		"--", "----", "--------", "----", "-----------")
	for _, job := range jobs.Activate {
		con.Printf(
			"%-4d %-16s %-10s %-6d %s\n",
			job.ID,
			job.Name,
			job.Protocol,
			job.Port,
			job.Description,
		)
	}

}
