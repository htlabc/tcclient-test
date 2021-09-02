package temple

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test_newCyclebarrier(t *testing.T) {

	bafun := func() error {
		fmt.Println("already execute 5 times,exec next...")
		return nil
	}

	stopch := make(chan struct{})

	go func() {
		time.Sleep(10 * time.Second)
		stopch <- struct{}{}
	}()

	cycle := newCyclebarrier(5, bafun, stopch)
	for i := 0; i < 10; i++ {
		args := make([]interface{}, 0)
		args = append(args, i)
		execfunc := func(args ...interface{}) (interface{}, error) {
			//fmt.Println(reflect.TypeOf(args[0].([]interface{})[0]).Kind())
			fmt.Println("execfunc :" + strconv.Itoa(args[0].([]interface{})[0].(int)))
			return "name" + strconv.Itoa(args[0].([]interface{})[0].(int)), nil
		}
		cycle.WithArgsAndFuncs(args, execfunc)
	}

	result := cycle.CyclicBarrier()

	fmt.Println(result.Result)

}
