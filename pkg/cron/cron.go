package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskInfo struct {
	Name        string
	StartTime   time.Time
	LastRunTime time.Time
	NextRunTime time.Time
	IsRunning   bool
}

type CronManager struct {
	cronInstance *cron.Cron
	tasks        map[cron.EntryID]*TaskInfo
	mu           sync.Mutex
}

func NewCronManager() *CronManager {
	return &CronManager{
		cronInstance: cron.New(),
		tasks:        make(map[cron.EntryID]*TaskInfo),
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

	id, err := cm.cronInstance.AddFunc(spec, func() {
		cm.mu.Lock()
		task := cm.tasks[id]
		task.IsRunning = true
		task.LastRunTime = time.Now()
		cm.mu.Unlock()

		cmd()

		cm.mu.Lock()
		task.IsRunning = false
		task.NextRunTime = cm.cronInstance.Entry(id).Next
		cm.mu.Unlock()
	})
	if err != nil {
		return 0, err
	}

	entry := cm.cronInstance.Entry(id)
	cm.tasks[id] = &TaskInfo{
		Name:        taskName,
		StartTime:   time.Now(),
		NextRunTime: entry.Next,
	}
	return id, nil
}

func (cm *CronManager) RemoveTask(id cron.EntryID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)
}

func (cm *CronManager) ListTasks() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, task := range cm.tasks {
		fmt.Printf("Task ID: %d, Name: %s, StartTime: %v, LastRunTime: %v, NextRunTime: %v, IsRunning: %t\n",
			id, task.Name, task.StartTime, task.LastRunTime, task.NextRunTime, task.IsRunning)
	}
}

func (cm *CronManager) ModifyTask(id cron.EntryID, newSpec string, newCmd func()) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Remove the old task
	cm.cronInstance.Remove(id)
	delete(cm.tasks, id)

	// Add the new task
	newID, err := cm.cronInstance.AddFunc(newSpec, newCmd)
	if err != nil {
		return err
	}
	cm.tasks[newID] = cm.tasks[id]
	return nil
}
