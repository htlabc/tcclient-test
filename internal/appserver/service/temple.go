package service

import (
	"githup.com/htl/tcclienttest/internal/appserver/temple/temple"
	temple2 "githup.com/htl/tcclienttest/internal/pkg/temple"
)

type StrageTemple struct {
	temple.StrategyTempleFactoryImpl
}

func (s *StrageTemple) SetStrageTemple(impl temple.StrategyTempleFactoryImpl) {
	s.Specific = impl
}

func (s *StrageTemple) Run() (results *temple2.Results) {
	results = s.SchedGroup()
	results = s.Semaphore()
	results = s.CyclicBarrier()
	results = s.GroupAsyncWait()
	results = s.GroupSyncWait()

	return results
}
