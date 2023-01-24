package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/AndrOGennad/go-yandex-devops/internal/server"
)

func main() {
	parentCtx := context.Background()
	ctx, cnl := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cnl()

	_ = server.Run(ctx)
	<-ctx.Done()
}
