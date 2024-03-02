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

func TestNatsPublisher(t *testing.T) {
	t.Run(".Name", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			js := mockjetstream.NewJetStream(t)
			publisher := &NatsPublisher{
				name:   testdata.Fake.App().Name(),
				conf:   testconf("nats://127.0.0.1:42222"),
				logger: testify.Logger(),
				js:     js,
			}
			require.Equal(sst, publisher.name, publisher.Name())
		})
	})

	t.Run(".Pub", func(st *testing.T) {
		st.Run("KO - validation error", func(sst *testing.T) {
			js := mockjetstream.NewJetStream(sst)
			publisher := &NatsPublisher{
				name:   testdata.Fake.App().Name(),
				conf:   testconf("nats://127.0.0.1:42222"),
				logger: testify.Logger(),
				js:     js,
			}

			ctx := context.Background()

			data := testdata.NewUser(clock.New())
			id := uuid.NewString()
			events := map[string]*entities.Event{
				id: {
					Subject:  testdata.Fake.Internet().Email(),
					Id:       id,
					Data:     data.Bytes(),
					Metadata: map[string]string{},
				},
			}

			errs := publisher.Pub(ctx, events)
			require.Equal(sst, len(events), len(errs))

			require.ErrorContains(sst, errs[id], "STREAMING.ENTITIES.EVENT")
		})

		st.Run("KO - publish error", func(sst *testing.T) {
			js := mockjetstream.NewJetStream(sst)
			publisher := &NatsPublisher{
				name:   testdata.Fake.App().Name(),
				conf:   testconf("nats://127.0.0.1:42222"),
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

			msg := NatsMsgFromEvent(events[id].Subject, events[id])
			js.EXPECT().PublishMsg(mock.Anything, msg).Return(nil, testdata.ErrorGeneric).Once()

			ctx := context.Background()
			errs := publisher.Pub(ctx, events)
			require.Equal(sst, len(events), len(errs))

			require.ErrorIs(sst, errs[id], testdata.ErrorGeneric)
		})

		st.Run("KO - duplicated error", func(sst *testing.T) {
			js := mockjetstream.NewJetStream(sst)
			publisher := &NatsPublisher{
				name:   testdata.Fake.App().Name(),
				conf:   testconf("nats://127.0.0.1:42222"),
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

			msg := NatsMsgFromEvent(events[id].Subject, events[id])
			ack := &jetstream.PubAck{
				Stream:    streamname(),
				Sequence:  testdata.Fake.UInt64(),
				Duplicate: true,
			}
			js.EXPECT().PublishMsg(mock.Anything, msg).Return(ack, nil).Once()

			ctx := context.Background()
			errs := publisher.Pub(ctx, events)
			require.Equal(sst, len(events), len(errs))

			require.ErrorContains(sst, errs[id], "STREAMING.PUBLISHER.EVENT_DUPLICATED")
		})
	})
}
