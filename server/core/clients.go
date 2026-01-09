package core

import (
	"sync"

	"github.com/longzekun/specter/client/constants"
	"github.com/longzekun/specter/protobuf/clientpb"
)

var (
	Clients = &clients{
		mu:     &sync.Mutex{},
		active: map[string]*Client{},
	}
)

type Client struct {
	ID       string
	Operator *clientpb.Operator
}

func (c *Client) ToProtoBuf() *clientpb.Client {
	return &clientpb.Client{
		ID:       c.ID,
		Name:     c.Operator.Name,
		Operator: c.Operator,
	}
}

func NewClient(operatorName string, uuid string) *Client {
	return &Client{
		ID: uuid,
		Operator: &clientpb.Operator{
			Name: operatorName,
		},
	}
}

type clients struct {
	mu     *sync.Mutex
	active map[string]*Client
}

func (c *clients) Add(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.active[client.ID] = client

	//	publish client joined message
	EventBroker.Publish(Event{
		EventType: constants.ClientJoinType,
		Client:    client,
	})
}

func (c *clients) Remove(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.active, client.ID)

	//	publish client quit message
	EventBroker.Publish(Event{
		EventType: constants.ClientLeaveType,
		Client:    client,
	})
}

func (c *clients) GetAllOperator() []*clientpb.Operator {
	var operators []*clientpb.Operator
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, client := range c.active {
		operators = append(operators, client.Operator)
	}
	return operators
}
