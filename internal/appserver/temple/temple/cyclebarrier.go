package temple

import (
	"context"
	"github.com/marusama/cyclicbarrier"
	"githup.com/htl/tcclienttest/internal/pkg/temple"
	"githup.com/htl/tcclienttest/pkg/log"
	"sync"
)

type cyclebarrierTemple struct {
	StrategyTempleFactoryImpl
	cyclicbarrier cyclicbarrier.CyclicBarrier
	execFunc      []argsAndFuncs
	sync.Mutex
	stopch chan struct{}
}

func (c *cyclebarrierTemple) WithArgsAndFuncs(args []interface{}, fc func(args ...interface{}) (interface{}, error)) {
	c.execFunc = append(c.execFunc, argsAndFuncs{args: args, fc: fc})
}

type argsAndFuncs struct {
	args interface{}
	fc   func(args ...interface{}) (interface{}, error)
}

func newCyclebarrier(parties int, barrierAction func() error, stopch chan struct{}) *cyclebarrierTemple {
	return &cyclebarrierTemple{cyclicbarrier: cyclicbarrier.NewWithAction(parties, barrierAction), execFunc: make([]argsAndFuncs, 0), stopch: stopch}
}

func (c *cyclebarrierTemple) CyclicBarrier() (result *temple.Results) {
	result = &temple.Results{make([]interface{}, 0), make([]error, 0)}
	for _, fn := range c.execFunc {
		go func(funcs argsAndFuncs) {
			data, err := funcs.fc(funcs.args)
			result.Errors = append(result.Errors, err)
			result.Result = append(result.Result, data)
			log.Info("exec cyclebarrier temple...")
			c.cyclicbarrier.Await(context.TODO())
		}(fn)
	}
	select {
	case <-c.stopch:
		log.Info("stop cyclebarrier task...")
		return result
	}
}
