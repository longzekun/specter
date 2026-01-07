package rpc

import (
	"context"

	"github.com/longzekun/specter/client/constants"
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/server/c2"
)

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
