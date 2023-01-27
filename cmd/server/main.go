package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/AndrOGennad/go-yandex-devops/internal/server"
)

func main() {
	parentCtx := context.Background()
	ctx, cnl := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cnl()

	if err := server.Run(ctx); err != nil {
		fmt.Println(err)
		return
	}
	<-ctx.Done()
	fmt.Println("сервер остановлен")
}
