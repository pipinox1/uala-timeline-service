package main

import (
	"fmt"
	httpServer "net/http"
	"os"
	"os/signal"
	"syscall"
	"uala-timeline-service/cmd/consumer"
	"uala-timeline-service/cmd/http"
	"uala-timeline-service/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config file: %w", err))
	}

	dependencies, err := config.BuildDependencies(*cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error building dependencies: %w", err))
	}

	c, subscriptions := consumer.SetupConsumer(cfg, dependencies)
	router := http.SetupRouterAndRoutes(cfg, dependencies)
	go func() {
		fmt.Println("starting server on port: ", cfg.Port)
		if err := httpServer.ListenAndServe(
			fmt.Sprintf(":%s", cfg.Port),
			router,
		); err != nil {
			fmt.Println("error starting server")
			panic(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	c.Close()
	for _, s := range subscriptions {
		_ = s.Unsubscribe()
	}

}
