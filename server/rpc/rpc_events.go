package rpc

import (
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/core"
)

func (s *Server) Events(_ *commonpb.Empty, stream rpcpb.SpecterRPC_EventsServer) error {
	events := core.EventBroker.Subscribe()

	defer func() {
		core.EventBroker.Unsubscribe(events)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-events:
			pbEvent := &clientpb.Event{
				EventType: event.EventType,
			}

			err := stream.Send(pbEvent)
			if err != nil {
				return err
			}
		}
	}
}
