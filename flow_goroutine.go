/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2022-11-02 13:39:06
# File Name: flow_goroutine.go
# Description:
####################################################################### */

package flow

import (
	"fmt"
	"sync"
)

type Goroutine interface {
	Run(ctx FlowContext) (r int, err error)
}

type GoroutineFunc func(ctx FlowContext) (r int, err error)

func (this GoroutineFunc) Run(ctx FlowContext) (r int, err error) { return this(ctx) }

type GoroutineFlow struct {
	*CommonFlow
	flowGoroutine Goroutine
}

func NewGoroutineFlow() *GoroutineFlow {
	o := &GoroutineFlow{}
	o.CommonFlow = NewCommonFlow()
	return o
}

func (this *GoroutineFlow) AddPloy(ploys ...Ploy) *GoroutineFlow {
	for _, ploy := range ploys {
		this.CommonFlow.AddPloy(ploy)
	}
	return this
}

func (this *GoroutineFlow) AddPloyFunc(fns ...func(ctx FlowContext)) *GoroutineFlow {
	for _, fn := range fns {
		this.AddPloy(PloyFunc(fn))
	}
	return this
}

func (this *GoroutineFlow) AddFlow(flows ...Flow) *GoroutineFlow {
	for _, flow := range flows {
		this.CommonFlow.AddFlow(flow)
	}
	return this
}

func (this *GoroutineFlow) AddHook(hooks ...Hook) *GoroutineFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *GoroutineFlow) SetGoroutine(flowGoroutine Goroutine) (r *GoroutineFlow) {
	this.flowGoroutine = flowGoroutine
	return this
}
func (this *GoroutineFlow) SetGoroutineFunc(fn func(ctx FlowContext) (r int, err error)) (r *GoroutineFlow) {
	this.SetGoroutine(GoroutineFunc(fn))
	return this
}

func (this *GoroutineFlow) Run(ctx FlowContext) {
	if this.flowGoroutine == nil {
		ctx.AddError(fmt.Errorf("goroutine flow must set goroutine first"))
		return
	}

	processNum, err := this.flowGoroutine.Run(ctx)
	if err != nil {
		ctx.AddError(fmt.Errorf("goroutine cond call error: %s", err))
		return
	}

	var wg sync.WaitGroup
	for ; processNum > 0; processNum-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, runnable := range this.runnables {
				if ctx.HasError() {
					break
				}
				this.hookMgr.before(ctx, runnable)
				runnable.Run(ctx)
				this.hookMgr.after(ctx, runnable)
			}
		}()
	}
	wg.Wait()
}
