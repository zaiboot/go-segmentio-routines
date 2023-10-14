package kafka

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
	/*fetching message: failed to dial: failed to open connection to localhost:9093: dial tcp [::1]:9093:
	connectex: No connection could be made because the target machine actively refused it
	*/

	wait_time := 1 * time.Second
	kafka_unavailable_count := 0
	for {
		err := ctx.Err()
		if err != nil {
			subLogger.Fatal().
				Str("error", err.Error()).
				Msg("Closed due to other process request")
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			subLogger.Err(err).
				Msg("Unable to read a message")
			if strings.Contains(err.Error(), "failed to dial") {
				// it looks like kafka is down, we can try to wait for a `wait_time`
				kafka_unavailable_count += 1
				time.Sleep(wait_time * time.Duration(kafka_unavailable_count))
				subLogger.Debug().
					Int("kafka.down.count", kafka_unavailable_count).
					Msg("Backing off when kafka is down")
			}
			continue
		}

		kafka_unavailable_count = 0

		err = f(string(msg.Value), &subLogger)
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
