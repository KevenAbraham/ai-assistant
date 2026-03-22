package main

import (
	"go.uber.org/fx"

	"github.com/KevenAbraham/ai-assistant/cmd/api/modules"
)

func main() {
	fx.New(
		modules.ConfigModule,
		modules.DBModule,
		modules.AIModule,
		modules.ServerModule,
	).Run()
}
