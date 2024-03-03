package streaming

import (
	"strings"

	"github.com/kanthorlabs/common/streaming/entities"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NatsMsgFromEvent(event *entities.Event) *natsio.Msg {
	msg := &natsio.Msg{
		Subject: event.Subject,
		Header: natsio.Header{
			natsio.MsgIdHdr: []string{event.Id},
		},
		Data: event.Data,
	}
	for key, value := range event.Metadata {
		msg.Header.Set(key, value)
	}

	return msg
}
func NatsMsgToEvent(msg jetstream.Msg) *entities.Event {
	event := &entities.Event{
		Subject:  msg.Subject(),
		Id:       msg.Headers().Get(natsio.MsgIdHdr),
		Data:     msg.Data(),
		Metadata: map[string]string{},
	}
	for key, value := range msg.Headers() {
		if strings.HasPrefix(key, "Nats") {
			continue
		}
		event.Metadata[key] = value[0]
	}
	return event
}
