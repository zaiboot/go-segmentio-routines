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
	dedup := &custom_logic.Redis_dedupe{}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go infra.StartReadynessWebServer(config, &l, &wg, ctx, cancel)
	wg.Add(1)
	go infra.MonitorProcesses(&l, &wg, ctx, cancel)

	wg.Add(1)
	go kafka.Consume(&wg, &l, ctx, cancel, dedup, config, custom_logic.DoWork)

	wg.Wait()
	cancel()
}
