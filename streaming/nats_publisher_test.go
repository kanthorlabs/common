package streaming

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/clock"
	mockjetstream "github.com/kanthorlabs/common/mocks/jetstream"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNatsPublisher_Name(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(t)
		publisher := &NatsPublisher{
			name:   pubsubname(),
			conf:   testconf(testdata.NatsUri),
			logger: testify.Logger(),
			js:     js,
		}
		require.Equal(st, publisher.name, publisher.Name())
	})
}

func TestNatsPublisher_Pub(t *testing.T) {
	t.Run("KO - validation error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		publisher := &NatsPublisher{
			name:   pubsubname(),
			conf:   testconf(testdata.NatsUri),
			logger: testify.Logger(),
			js:     js,
		}

		ctx := context.Background()

		data := testdata.NewUser(clock.New())
		id := uuid.NewString()
		events := map[string]*entities.Event{
			id: {
				Subject:  subjectname() + "+" + subjectname(),
				Id:       id,
				Data:     data.Bytes(),
				Metadata: map[string]string{},
			},
		}

		errs := publisher.Pub(ctx, events)
		require.Equal(st, len(events), len(errs))

		require.ErrorContains(st, errs[id], "STREAMING.EVENT")
	})

	t.Run("KO - publish error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		publisher := &NatsPublisher{
			name:   pubsubname(),
			conf:   testconf(testdata.NatsUri),
			logger: testify.Logger(),
			js:     js,
		}

		data := testdata.NewUser(clock.New())
		id := uuid.NewString()
		events := map[string]*entities.Event{
			id: {
				Subject:  subjectname(),
				Id:       id,
				Data:     data.Bytes(),
				Metadata: map[string]string{},
			},
		}

		msg := NatsMsgFromEvent(events[id])
		js.EXPECT().PublishMsg(mock.Anything, msg).Return(nil, testdata.ErrGeneric).Once()

		ctx := context.Background()
		errs := publisher.Pub(ctx, events)
		require.Equal(st, len(events), len(errs))

		require.ErrorIs(st, errs[id], testdata.ErrGeneric)
	})

	t.Run("KO - duplicated error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		publisher := &NatsPublisher{
			name:   pubsubname(),
			conf:   testconf(testdata.NatsUri),
			logger: testify.Logger(),
			js:     js,
		}

		data := testdata.NewUser(clock.New())
		id := uuid.NewString()
		events := map[string]*entities.Event{
			id: {
				Subject:  subjectname(),
				Id:       id,
				Data:     data.Bytes(),
				Metadata: map[string]string{},
			},
		}

		msg := NatsMsgFromEvent(events[id])
		ack := &jetstream.PubAck{
			Stream:    streamname(),
			Sequence:  testdata.Fake.UInt64(),
			Duplicate: true,
		}
		js.EXPECT().PublishMsg(mock.Anything, msg).Return(ack, nil).Once()

		ctx := context.Background()
		errs := publisher.Pub(ctx, events)
		require.Equal(st, len(events), len(errs))

		require.ErrorContains(st, errs[id], "STREAMING.PUBLISHER.EVENT_DUPLICATED")
	})
}
