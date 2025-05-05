package crons

import (
	"Taurus/pkg/cron"
	"Taurus/pkg/loggerx"
)

func init() {
	cron.CronManagerInstance.AddTask("*/2 * * * * *", "DemoCron", func() {
		loggerx.DefaultLogger.Info("demo crond been executed")
	})
}
