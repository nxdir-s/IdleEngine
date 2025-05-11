package valobj

import "context"

type Event struct {
	Ctx      context.Context
	Consumed chan struct{}
	Msg      *Message
}

type Message struct {
	Value string `json:"value"`
}
