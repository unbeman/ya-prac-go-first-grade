package worker

type workersPool struct {
}

func NewWorkersPool() (*workersPool, error) {
	return &workersPool{}, nil
}
