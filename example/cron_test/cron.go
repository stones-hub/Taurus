package main

import (
	"Taurus/pkg/cron"
	"Taurus/pkg/util"
	"log"
)

// 1.秒（可选）: 0-59
// 2.分钟: 0-59
// 3.小时: 0-23
// 4.日: 1-31
// 5.月: 1-12 或者 JAN-DEC
// 6.星期几: 0-6 (0 表示周日) 或者 SUN-SAT

func main() {
	// 每分钟执行一次: * * * * *
	// 每小时的第 5 分钟执行: 5 * * * *
	// 每天的凌晨 0 点执行: 0 0 * * *
	// 每周一的凌晨 0 点执行: 0 0 * * 1
	// 每个月的 1 号凌晨 0 点执行: 0 0 1 * *
	// 每隔 5 分钟执行一次: */5 * * * *
	// 每隔 10 秒执行一次（需要启用秒字段支持）: */10 * * * * *
	cron.CronManagerInstance.AddTask("*/10 * * * * *", "每隔 10 秒执行一次", func() {
		log.Println("每隔 10 秒执行一次")
	})

	// 每分钟的第 5 秒执行
	cron.CronManagerInstance.AddTask("5 * * * * *", "每分钟的第 5 秒执行", func() {
		log.Println("每分钟的第 5 秒执行")

		// 获取任务所有信息
		taskStatuses := cron.CronManagerInstance.ListTasks()

		// 准备表格数据
		headers := []string{"任务ID", "任务名称", "开始时间", "上次运行时间", "下次运行时间"}
		lines := make([][]interface{}, 0)
		for _, taskStatus := range taskStatuses {
			lines = append(lines, []interface{}{
				taskStatus.ID,
				taskStatus.Name,
				taskStatus.StartTime.Format("2006-01-02 15:04:05"),
				taskStatus.PrevRun.Format("2006-01-02 15:04:05"),
				taskStatus.NextRun.Format("2006-01-02 15:04:05"),
			})
		}

		// 渲染表格
		util.RenderTable(headers, lines)
	})

	// 启动 cron 调度器
	cron.CronManagerInstance.Start()

	// 阻止主线程退出
	select {}
}
