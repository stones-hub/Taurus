package crons

import (
	"Taurus/pkg/cron"
	"Taurus/pkg/logx"
)

func init() {
	cron.Core.AddTask("*/2 * * * * *", "DemoCron", func() {
		logx.Core.Info("default", "demo crond been executed")
	})
}
