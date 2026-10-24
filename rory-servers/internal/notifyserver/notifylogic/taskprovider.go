package notifylogic

import "example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic/task"

// TaskProvider ?
type TaskProvider interface {
	GetTask() task.Task
}

// ExternalTaskPrvider 外網工作任務
type ExternalTaskPrvider struct{}

// GetTask ?
func (etp *ExternalTaskPrvider) GetTask() task.Task {
	return nil
}

// InternalTaskProvider 內網工作任務
type InternalTaskProvider struct{}

// GetTask ?
func (itp *InternalTaskProvider) GetTask() task.Task {
	tg := new(task.Group)
	tg.AddTask(task.NewAnnTask(5, 1))
	tg.AddTask(task.NewAttrTask(5, 1))
	//tg.AddTask(task.NewRefreshTask(300, 1))
	tg.AddTask(task.NewCmdWorkerTask(5))
	return tg
}
