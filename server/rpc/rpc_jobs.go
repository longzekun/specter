package rpc

import (
	"context"

	"github.com/longzekun/specter/client/constants"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/server/c2"
	"github.com/longzekun/specter/server/core"
)

func (s *Server) GetAllJobs(ctx context.Context, _ *commonpb.Empty) (*clientpb.Jobs, error) {
	//	get all jobs
	retJobs := &clientpb.Jobs{}
	jobs := core.Jobs.All()
	for _, job := range jobs {
		activateJob := clientpb.Job{
			ID:          int32(job.ID),
			Name:        job.Name,
			Description: job.Description,
			Protocol:    job.Protocol,
			Port:        int32(job.Port),
		}
		retJobs.Activate = append(retJobs.Activate, &activateJob)
	}

	return retJobs, nil
}

func (s *Server) KillAllJobs(_ context.Context, req *commonpb.Empty) (*commonpb.Empty, error) {
	jobs := core.Jobs.All()
	for _, job := range jobs {
		job.JobCtl <- true
	}
	return nil, nil
}

func (s *Server) KillJob(_ context.Context, killJob *clientpb.KillJobReq) (*commonpb.Empty, error) {
	for _, job := range core.Jobs.All() {
		if job.ID == int(killJob.ID) {
			job.JobCtl <- true
		}
	}
	return nil, nil
}

func (s *Server) StartMTLSListener(ctx context.Context, req *clientpb.MTLSListenerReq) (*clientpb.ListenerJob, error) {
	job, err := c2.StartMTLSListenerJob(req.Host, req.Port)
	if err != nil {
		return nil, err
	}

	listenerJob := &clientpb.ListenerJob{
		JobID:    uint32(job.ID),
		Type:     constants.MtlsStr,
		MTLSConf: req,
	}

	return listenerJob, nil
}
