package kafka

import (
	"context"
	"strings"
	"sync"

	"zaiboot/segmentIO.tests/internal/configs"
	"zaiboot/segmentIO.tests/internal/custom_logic"

	segmentio "github.com/segmentio/kafka-go"

	"github.com/rs/zerolog"
)

func getKafkaReader(brokers, topic, groupID string, l *zerolog.Logger) *segmentio.Reader {
	return segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{brokers},
		GroupID:     groupID,
		Topic:       topic,
		Logger:      l,
		ErrorLogger: l,
	})
}
func Consume(wg *sync.WaitGroup, l *zerolog.Logger, ctx context.Context, cancel context.CancelFunc, d custom_logic.Dedupe_storage, config *configs.Config,
	f func(ctx context.Context, message string, l *zerolog.Logger, d custom_logic.Dedupe_storage, config *configs.Config) error) {
	defer wg.Done()

	subLogger := l.With().
		Str("bootstrap.servers", config.KAFKA_BROKER).
		Str("topic", config.KAFKA_TOPIC).
		Str("groupdId", config.KAFKA_GROUPID).
		Logger()
	subLogger.Info().Msg("Start")
	reader := getKafkaReader(config.KAFKA_BROKER, config.KAFKA_TOPIC, config.KAFKA_GROUPID, l)
	defer reader.Close()

	for {
		err := ctx.Err()
		if err != nil {
			subLogger.Fatal().
				Str("error", err.Error()).
				Msg("Closed due to other process request")
		}

		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			subLogger.Err(err).
				Msg("Unable to read a message")
			if strings.Contains(err.Error(), "failed to dial") {
				// it looks like kafka is down, we can try to wait for a `wait_time`
				subLogger.Debug().
					Msg("Backing off when kafka is down")
			}
			continue
		}

		err = f(ctx, string(msg.Value), &subLogger, d, config)
		if err != nil {
			subLogger.Fatal().
				Str("error", err.Error()).
				Msg("Unable to process the message")

			break
		}

		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			subLogger.Err(err).
				Msg("Unable to commit")
		}

	}
}
