/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 15:15:27
# File Name: flow_pipe.go
# Description:
####################################################################### */

package flow

type PipeFlow struct {
	*CommonFlow
}

func NewPipeFlow() *PipeFlow {
	o := &PipeFlow{}
	o.CommonFlow = NewCommonFlow()
	return o
}

func (this *PipeFlow) AddPloy(ploy Ploy) *PipeFlow {
	this.CommonFlow.AddPloy(ploy)
	return this
}

func (this *PipeFlow) AddPloyFunc(fn func(ctx FlowContext)) *PipeFlow {
	this.AddPloy(PloyFunc(fn))
	return this
}

func (this *PipeFlow) AddFlow(flow Flow) *PipeFlow {
	this.CommonFlow.AddFlow(flow)
	return this
}

func (this *PipeFlow) AddHook(hooks ...Hook) *PipeFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *PipeFlow) Run(ctx FlowContext) {
	if ctx.HasError() {
		return
	}

	for _, runnable := range this.runnables {
		this.hookMgr.before(ctx, runnable)
		runnable.Run(ctx)
		this.hookMgr.after(ctx, runnable)
	}

	return
}