package main

import (
	"context"
	"os"
	"sync"
	"zaiboot/segmentIO.tests/internal/configs"
	"zaiboot/segmentIO.tests/internal/custom_logic"
	"zaiboot/segmentIO.tests/internal/infra"
	"zaiboot/segmentIO.tests/internal/kafka"
	"zaiboot/segmentIO.tests/internal/logger"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	l := logger.InitLog(config.APPLICATIONNAME, config.LOGLEVEL)

	l.Debug().Int("pid", os.Getpid()).Msg("Pid")

	var wg sync.WaitGroup
	message := make(chan string, 3)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go infra.StartReadynessWebServer(config, &l, &wg, ctx, cancel)
	wg.Add(1)
	go infra.MonitorProcesses(&l, &wg, ctx, cancel)
	kc, err := config.GetKafkaConfig()
	if err != nil {
		l.Fatal().Msg("Error getting kafka config")
	}
	c, err := kafka.BuildConsumer(kc[0], &l)
	if err != nil {
		l.Fatal().Msg("Unable to create consumer")
	}

	for _, k := range kc {
		wg.Add(1)
		go kafka.Consume(&wg, c, &l, ctx, cancel, k, custom_logic.DoWork)
	}

	wg.Wait()
	c.Close()
	close(message)
}