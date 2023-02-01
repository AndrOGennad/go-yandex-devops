package main

import (
	"context"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndrOGennad/go-yandex-devops/internal/agent"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	parentCtx := context.Background()
	ctx, cnl := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cnl()

	_ = agent.Run(ctx)
	<-ctx.Done()
}
