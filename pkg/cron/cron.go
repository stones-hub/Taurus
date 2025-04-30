package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskStatus struct {
	IsRunning   bool
	LastRunTime time.Time
	NextRunTime time.Time
}

type CronManager struct {
	cronInstance *cron.Cron
	tasks        map[cron.EntryID]string
	status       map[cron.EntryID]*TaskStatus
	mu           sync.Mutex
}

func Initialize() *CronManager {
	return &CronManager{
		cronInstance: cron.New(),
		tasks:        make(map[cron.EntryID]string),
		status:       make(map[cron.EntryID]*TaskStatus),
	}
}

func (cm *CronManager) Start() {
	cm.cronInstance.Start()
}

func (cm *CronManager) Stop() {
	cm.cronInstance.Stop()
}

func (cm *CronManager) AddTask(spec string, taskName string, cmd func()) (cron.EntryID, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 添加任务并获取任务ID
	id, err := cm.cronInstance.AddFunc(spec, func() {
		cm.mu.Lock()
		taskStatus := cm.status[id]
		taskStatus.IsRunning = true
		taskStatus.LastRunTime = time.Now()
		cm.mu.Unlock()

		// 执行任务
		cmd()

		cm.mu.Lock()
		taskStatus.IsRunning = false
		taskStatus.NextRunTime = cm.cronInstance.Entry(id).Next
		cm.mu.Unlock()
	})
	if err != nil {
		return 0, err
	}

	// 获取任务的下次运行时间
	entry := cm.cronInstance.Entry(id)
	nextRunTime := entry.Next

	// 记录任务状态
	cm.tasks[id] = taskName
	cm.status[id] = &TaskStatus{IsRunning: false, NextRunTime: nextRunTime, LastRunTime: time.Now()}
	return id, nil
}

func (cm *CronManager) RemoveTask(id cron.EntryID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)
	delete(cm.status, id)
}

func (cm *CronManager) ListTasks() map[cron.EntryID]*TaskStatus {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, name := range cm.tasks {
		status := cm.status[id]
		fmt.Printf("Task ID: %d, Name: %s, IsRunning: %t, LastRunTime: %v, NextRunTime: %v\n",
			id, name, status.IsRunning, status.LastRunTime, status.NextRunTime)
	}

	return cm.status
}

func (cm *CronManager) ModifyTask(id cron.EntryID, newSpec string, newCmd func()) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 获取原有的任务名称
	taskName, exists := cm.tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d does not exist", id)
	}

	// Remove the old task
	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)
	delete(cm.status, id)

	// Add the new task
	newID, err := cm.cronInstance.AddFunc(newSpec, newCmd)
	if err != nil {
		return err
	}

	// 获取新任务的下次运行时间
	entry := cm.cronInstance.Entry(newID)
	nextRunTime := entry.Next

	// 保留原有的任务名称并记录新任务状态
	cm.tasks[newID] = taskName
	cm.status[newID] = &TaskStatus{IsRunning: false, NextRunTime: nextRunTime, LastRunTime: time.Now()}
	return nil
}

var (
	DefaultCronManager *CronManager
)

func init() {
	DefaultCronManager = Initialize()
}
