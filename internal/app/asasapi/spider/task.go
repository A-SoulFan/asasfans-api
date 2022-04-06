package spider

type Task interface {
	Run()
	Stop() error
}

type TaskManager struct {
	stopChan chan bool
	taskList []Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		stopChan: make(chan bool),
		taskList: make([]Task, 0, 5),
	}
}
