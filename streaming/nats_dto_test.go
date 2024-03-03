package streaming

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/mocks/jetstream"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/testdata"
	natsio "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestNats_DataTransferObject(t *testing.T) {
	id := uuid.NewString()
	in := &entities.Event{
		Subject: subjectname(),
		Id:      id,
		Data:    []byte("data"),
		Metadata: map[string]string{
			"User-Agent": testdata.Fake.UserAgent().InternetExplorer(),
		},
	}

	msg := NatsMsgFromEvent(in)
	require.Equal(t, in.Subject, msg.Subject)
	require.Equal(t, in.Id, msg.Header.Get(natsio.MsgIdHdr))

	out := NatsMsgToEvent(mockjsmsg(t, in))
	require.Equal(t, in.String(), out.String())
}

func mockjsmsg(t *testing.T, in *entities.Event) *jetstream.Msg {
	jsmsg := jetstream.NewMsg(t)
	jsmsg.EXPECT().Subject().Return(in.Subject).Times(1)
	jsmsg.EXPECT().Headers().Return(natsio.Header{
		natsio.MsgIdHdr: []string{in.Id},
		"User-Agent":    []string{in.Metadata["User-Agent"]},
	}).Times(2)
	jsmsg.EXPECT().Data().Return(in.Data).Times(1)

	return jsmsg
}
