/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 17:51:19
# File Name: flow_parallel.go
# Description:
####################################################################### */

package flow

import (
	"sync"
)

type ParallelFlow struct {
	*CommonFlow
	maxProcessNum int
}

func NewParallelFlow() *ParallelFlow {
	o := &ParallelFlow{}
	o.CommonFlow = NewCommonFlow()
	o.maxProcessNum = 10
	return o
}

func (this *ParallelFlow) AddPloy(ploys ...Ploy) *ParallelFlow {
	for _, ploy := range ploys {
		this.CommonFlow.AddPloy(ploy)
	}
	return this
}

func (this *ParallelFlow) AddPloyFunc(fns ...func(ctx FlowContext)) *ParallelFlow {
	for _, fn := range fns {
		this.AddPloy(PloyFunc(fn))
	}
	return this
}

func (this *ParallelFlow) AddFlow(flows ...Flow) *ParallelFlow {
	for _, flow := range flows {
		this.CommonFlow.AddFlow(flow)
	}
	return this
}

func (this *ParallelFlow) AddHook(hooks ...Hook) *ParallelFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *ParallelFlow) SetMaxProcess(num int) (r *ParallelFlow) {
	this.maxProcessNum = num
	return this
}

func (this *ParallelFlow) Run(ctx FlowContext) {
	if ctx.HasError() {
		return
	}

	fn := func(ctx FlowContext, runnableChan chan Runnable) {
		for {
			select {
			default:
				return
			case runnable := <-runnableChan:
				this.hookMgr.before(ctx, runnable)
				runnable.Run(ctx)
				this.hookMgr.after(ctx, runnable)
			}
		}
	}

	processNum := this.maxProcessNum
	if this.maxProcessNum > len(this.runnables) {
		processNum = len(this.runnables)
	}

	runnableChan := make(chan Runnable, len(this.runnables))
	for _, runnable := range this.runnables {
		runnableChan <- runnable
	}

	var wg sync.WaitGroup
	for i := 0; i < processNum; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); fn(ctx, runnableChan) }()
	}
	wg.Wait()

	return
}
