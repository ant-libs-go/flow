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

func (this *PipeFlow) AddPloy(ploys ...Ploy) *PipeFlow {
	for _, ploy := range ploys {
		this.CommonFlow.AddPloy(ploy)
	}
	return this
}

func (this *PipeFlow) AddPloyFunc(fns ...func(ctx FlowContext)) *PipeFlow {
	for _, fn := range fns {
		this.AddPloy(PloyFunc(fn))
	}
	return this
}

func (this *PipeFlow) AddFlow(flows ...Flow) *PipeFlow {
	for _, flow := range flows {
		this.CommonFlow.AddFlow(flow)
	}
	return this
}

func (this *PipeFlow) AddHook(hooks ...Hook) *PipeFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *PipeFlow) Run(ctx FlowContext) {
	for _, runnable := range this.runnables {
		if ctx.HasError() {
			break
		}
		this.hookMgr.before(ctx, runnable)
		runnable.Run(ctx)
		this.hookMgr.after(ctx, runnable)
	}
}
