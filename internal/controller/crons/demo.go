package crons

import (
	"Taurus/pkg/cron"
	"Taurus/pkg/loggerx"
)

func init() {
	cron.Core.AddTask("*/2 * * * * *", "DemoCron", func() {
		loggerx.Core.Info("demo crond been executed")
	})
}
