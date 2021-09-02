package temple

import (
	temple "githup.com/htl/tcclienttest/internal/pkg/temple"
)

type StrategyTempleFactory interface {
	GroupSyncWait() *temple.Results
	GroupAsyncWait() *temple.Results
	CyclicBarrier() *temple.Results
	Semaphore() *temple.Results
	SchedGroup() *temple.Results
}

//absolut

type StrategyTempleFactoryImpl struct {
	Specific StrategyTempleFactory
}

func (s StrategyTempleFactoryImpl) GroupSyncWait() *temple.Results {
	return s.Specific.GroupSyncWait()
}

func (s StrategyTempleFactoryImpl) GroupAsyncWait() *temple.Results {
	return s.GroupAsyncWait()
}

func (s StrategyTempleFactoryImpl) CyclicBarrier() *temple.Results {
	return s.CyclicBarrier()
}

func (s StrategyTempleFactoryImpl) Semaphore() *temple.Results {
	return s.Semaphore()
}

func (s StrategyTempleFactoryImpl) SchedGroup() *temple.Results {
	return s.SchedGroup()
}

func NewStragegyTempleImpl(Specific StrategyTempleFactory) *StrategyTempleFactoryImpl {
	return &StrategyTempleFactoryImpl{Specific: Specific}
}
