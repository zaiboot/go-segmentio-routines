package main

import (
	"context"
	"os"
	"sync"
	"zaiboot/segmentIO.tests/internal/configs"
	"zaiboot/segmentIO.tests/internal/infra"
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
	wg.Wait()
	close(message)
}
