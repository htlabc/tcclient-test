package task

var client TaskFactory

type TaskFactory interface {
	PreExecTask() error
	ExecTask() error
	AfterExecTask() error
	DeleteTask() error
	InstallTask() error
}

func Client() TaskFactory {
	return client
}

func SetClient(factory TaskFactory) {
	client = factory
}
