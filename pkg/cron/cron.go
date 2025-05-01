package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// TaskInfo 结构体用于存储任务的基本信息
type TaskInfo struct {
	Name      string
	StartTime time.Time
}

// CronManager 管理所有的定时任务
type CronManager struct {
	cronInstance *cron.Cron                 // cron 实例
	tasks        map[cron.EntryID]*TaskInfo // 存储任务信息的映射
	mu           sync.Mutex                 // 保护共享资源的互斥锁
}

// NewCronManager 创建一个新的 CronManager 实例
func NewCronManager() *CronManager {
	return &CronManager{
		cronInstance: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger), // 使用默认日志记录器恢复任务
		)),
		tasks: make(map[cron.EntryID]*TaskInfo),
		mu:    sync.Mutex{},
	}
}

// Start 启动 cron 调度器
func (cm *CronManager) Start() {
	cm.cronInstance.Start()
}

// Stop 停止 cron 调度器
func (cm *CronManager) Stop() {
	cm.cronInstance.Stop()
}

// AddTask 添加一个新的定时任务
func (cm *CronManager) AddTask(spec string, taskName string, cmd func()) (cron.EntryID, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id, err := cm.cronInstance.AddFunc(spec, cmd)
	if err != nil {
		return 0, err
	}

	// 创建任务信息
	cm.tasks[id] = &TaskInfo{
		Name:      taskName,
		StartTime: time.Now(),
	}
	return id, nil
}

// RemoveTask 移除一个定时任务
func (cm *CronManager) RemoveTask(id cron.EntryID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)
}

// ListTasks 列出所有的定时任务
func (cm *CronManager) ListTasks() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, task := range cm.tasks {
		entry := cm.cronInstance.Entry(id)
		fmt.Printf("Task ID: %d, Name: %s, StartTime: %v, NextRunTime: %v, PrevRunTime: %v\n",
			id, task.Name, task.StartTime, entry.Next, entry.Prev)
	}
}

// ModifyTask 修改一个定时任务
func (cm *CronManager) ModifyTask(id cron.EntryID, newSpec string, newCmd func()) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 获取任务信息
	taskInfo, exists := cm.tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d does not exist", id)
	}

	// Remove the old task
	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)

	_, err := cm.AddTask(newSpec, taskInfo.Name, newCmd)
	if err != nil {
		return err
	}
	return nil
}
