package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"zaiboot/segmentIO.tests/internal/configs"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

func BuildConsumer(kc configs.KafkaConfig, l *zerolog.Logger) (sarama.Consumer, error) {

	config := sarama.NewConfig()
	config.ClientID = "go-kafka-consumer"
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	//TODO: Use a setting
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	brokers := strings.Split(kc.Broker, ",")

	// Create new consumer
	c, err := sarama.NewConsumer(brokers, config)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func Consume(wg *sync.WaitGroup, c sarama.Consumer, l *zerolog.Logger, ctx context.Context, cancel context.CancelFunc, kc configs.KafkaConfig,
	f func(message string) error) {
	defer wg.Done()
	groupdId := fmt.Sprintf("%s-%d", kc.GroupId, kc.Partition)
	subLogger := l.With().
		Str("bootstrap.servers", kc.Broker).
		Str("topic", kc.Topic).
		Str("groupdId", groupdId).
		Int("partition", kc.Partition).Logger()
	// custom wpn consumer class to handle all
	consumer := Consumer{
		ready: make(chan bool),
	}

	saramaConfig := sarama.Config{}
	cg, err := sarama.NewConsumerGroup(strings.Split(kc.Broker, ","), groupdId, &saramaConfig)
	if err != nil {
		subLogger.Error().Err(err).Msg("Unable to consume partition")
		cancel()
		return
	}
	<-consumer.ready // Await till the consumer has been set up, can this be added
	for run := true; run; {

		if err := cg.Consume(ctx, strings.Split(kc.Topic, ","), &consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}

	}

	err = cg.Close()
	if err != nil {
		subLogger.Err(err).Msg("Unable to close the consumer")
	}

	l.Debug().Msg("Partition Consumer terminated")
}
