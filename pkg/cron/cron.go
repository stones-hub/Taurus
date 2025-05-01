package cron

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// TaskInfo 结构体用于存储任务的基本信息
type TaskInfo struct {
	Name      string
	Spec      string
	StartTime time.Time
}

// CronManager 管理所有的定时任务
type CronManager struct {
	cronInstance *cron.Cron                 // cron 实例
	tasks        map[cron.EntryID]*TaskInfo // 存储任务信息的映射
	mu           sync.Mutex                 // 保护共享资源的互斥锁
}

var (
	CronManagerInstance *CronManager
)

func init() {
	CronManagerInstance = InitializeCronManager()
}

// InitializeCronManager 创建一个新的 CronManager 实例
func InitializeCronManager() *CronManager {
	return &CronManager{
		cronInstance: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger), // 使用默认日志记录器恢复任务
		), cron.WithSeconds()),
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
		Spec:      spec,
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

type TaskStatus struct {
	ID        cron.EntryID // 任务ID
	Name      string       // 任务名称
	StartTime time.Time    // 任务开始时间
	PrevRun   time.Time    // 上次运行时间
	NextRun   time.Time    // 下次运行时间
}

// ListTasks 列出所有的定时任务
func (cm *CronManager) ListTasks() []*TaskStatus {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	taskStatuses := make([]*TaskStatus, 0)
	for id, task := range cm.tasks {
		entry := cm.cronInstance.Entry(id)
		taskStatus := &TaskStatus{
			ID:        id,
			Name:      task.Name,
			StartTime: task.StartTime,
			NextRun:   entry.Next,
			PrevRun:   entry.Prev,
		}
		taskStatuses = append(taskStatuses, taskStatus)
	}
	return taskStatuses
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

	if id, err := cm.AddTask(newSpec, taskInfo.Name, newCmd); err != nil {
		return err
	} else {
		log.Printf("Task %s modified successfully, New ID: %d", taskInfo.Name, id)
		return nil
	}
}
