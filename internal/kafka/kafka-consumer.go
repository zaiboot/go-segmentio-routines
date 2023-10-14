package kafka

import (
	"context"
	"fmt"
	"sync"

	"zaiboot/segmentIO.tests/internal/configs"

	segmentio "github.com/segmentio/kafka-go"

	"github.com/rs/zerolog"
)

func getKafkaReader(kafkaURL, topic, groupID string) *segmentio.Reader {
	return segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
func Consume(wg *sync.WaitGroup, l *zerolog.Logger, ctx context.Context, cancel context.CancelFunc, kc configs.KafkaConfig,
	f func(message string, l *zerolog.Logger) error) {
	defer wg.Done()
	groupdId := fmt.Sprintf("%s-%d", kc.GroupId, kc.Partition)
	subLogger := l.With().
		Str("bootstrap.servers", kc.Broker).
		Str("topic", kc.Topic).
		Str("groupdId", groupdId).
		Int("partition", kc.Partition).Logger()
	subLogger.Info().Msg("Start")
	reader := getKafkaReader(kc.Broker, kc.Topic, groupdId)
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			subLogger.Err(err).
				Msg("Unable to read a message")
		}
		err = f(string(msg.Value), &subLogger)
		if err != nil {
			subLogger.Fatal().
				Str("error", err.Error()).
				Msg("Unable to process the message")

			return
		}
		err = ctx.Err()
		if err != nil {
			subLogger.Fatal().
				Str("error", err.Error()).
				Msg("Closed due to other process request")
		}
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			subLogger.Err(err).
				Msg("Unable to commit")
		}

	}
}
