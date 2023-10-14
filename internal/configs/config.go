package configs

import (
	"strconv"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type KafkaConfig struct {
	Partition                int
	Broker                   string
	GroupId                  string
	Topic                    string
	BatchSize                int
	ConsumeTimeout           int
	MaxConsecutiveErrorCount int
}

type Config struct {
	PORT                                       string
	APPLICATIONNAME                            string
	LOGLEVEL                                   string
	KAFKA_PARTITIONS                           string
	KAFKA_BROKER                               string
	KAFKA_GROUPID                              string
	KAFKA_TOPIC                                string
	KAFKA_BATCH_SIZE                           int
	KAFKA_CONSUMER_DURATION_MS                 int
	KAFKA_CONSUMER_MAX_CONSECUTIVE_ERROR_COUNT int
}

// TODO: Define default values for kafka config.
func (c KafkaConfig) Init() KafkaConfig {
	c.MaxConsecutiveErrorCount = 10
	return c
}

// Can this be improved?
func (c *Config) GetKafkaConfig() ([]KafkaConfig, error) {
	kc := []KafkaConfig{}
	for _, p := range strings.Split(c.KAFKA_PARTITIONS, ",") {
		pnumer, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		kc = append(kc, KafkaConfig{
			Partition:                pnumer,
			Broker:                   c.KAFKA_BROKER,
			GroupId:                  c.KAFKA_GROUPID,
			Topic:                    c.KAFKA_TOPIC,
			BatchSize:                c.KAFKA_BATCH_SIZE,
			ConsumeTimeout:           c.KAFKA_CONSUMER_DURATION_MS,
			MaxConsecutiveErrorCount: c.KAFKA_CONSUMER_MAX_CONSECUTIVE_ERROR_COUNT,
		})

	}
	return kc, nil
}

func LoadConfig() (*Config, error) {
	var config Config
	var k = koanf.New(".")
	err := k.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil)
	if err != nil {
		return nil, err
	}
	err = k.Unmarshal("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
