package rpc

import (
	"github.com/longzekun/specter/protobuf/clientpb"
	"github.com/longzekun/specter/protobuf/commonpb"
	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/core"
	"github.com/longzekun/specter/server/db/models"
)

const (
	Transport = "transport"
	Operator  = "operator"
)

func (s *Server) Events(_ *commonpb.Empty, stream rpcpb.SpecterRPC_EventsServer) error {
	operator := stream.Context().Value(Operator).(*models.Operator)

	operatorClient := core.NewClient(operator.Name, operator.ID.String())
	core.Clients.Add(operatorClient)
	events := core.EventBroker.Subscribe()

	defer func() {
		core.EventBroker.Unsubscribe(events)
		core.Clients.Remove(operatorClient)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-events:
			pbEvent := &clientpb.Event{
				EventType: event.EventType,
			}

			if event.Client != nil {
				pbEvent.Client = event.Client.ToProtoBuf()
			}

			err := stream.Send(pbEvent)
			if err != nil {
				return err
			}
		}
	}
}
