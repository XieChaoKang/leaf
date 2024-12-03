package timer

import (
	"fmt"
	"leaf/conf"
	"leaf/log"
	"runtime"
	"sync"
	"time"
)

type dispatcherSet map[*Timer]struct{}

// one dispatcher per goroutine (goroutine not safe)
type Dispatcher struct {
	disMutex  sync.RWMutex
	ChanTimer chan *Timer
	closeFlag bool
}

func NewDispatcher(l int) *Dispatcher {
	disp := new(Dispatcher)
	disp.ChanTimer = make(chan *Timer, l)
	return disp
}

// Timer
type Timer struct {
	t    *time.Timer
	cb   interface{}
	args []interface{}
}

func (t *Timer) Stop() {
	t.t.Stop()
	t.cb = nil
	t.args = nil
}

func (t *Timer) Cb() interface{} {
	defer func() {
		t.cb = nil
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	if t.cb != nil {
		// execute
		switch t.cb.(type) {
		case func():
			t.cb.(func())()
		case func([]interface{}):
			t.cb.(func([]interface{}))(t.args)
		case func([]interface{}) interface{}:
			ret := t.cb.(func([]interface{}) interface{})(t.args)
			return ret
		case func([]interface{}) []interface{}:
			ret := t.cb.(func([]interface{}) []interface{})(t.args)
			return ret
		}
	}

	return nil
}

func (disp *Dispatcher) AfterFunc(d time.Duration, cb interface{}, args ...interface{}) *Timer {
	t := new(Timer)
	t.cb = cb
	t.args = args
	t.t = time.AfterFunc(d, func() {
		if !disp.closeFlag {
			TryCatch(func() {
				disp.ChanTimer <- t
			}, func(err interface{}) {
			})
		}
	})

	return t
}

func (disp *Dispatcher) OnClose() {
	disp.disMutex.Lock()
	defer disp.disMutex.Unlock()

	if !disp.closeFlag {
		close(disp.ChanTimer)
		disp.closeFlag = true
	}
}

// Cron
type Cron struct {
	t *Timer
}

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}

func (disp *Dispatcher) CronFunc(cronExpr *CronExpr, _cb func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer _cb()

		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = disp.AfterFunc(nextTime.Sub(now), cb)

	return c
}

func TryCatch(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%v %s", err, buf[:n])
			handler(stackInfo)
		}
	}()
	fun()
}
