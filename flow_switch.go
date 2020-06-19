/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 18:33:05
# File Name: flow_switch.go
# Description:
####################################################################### */

package flow

import "fmt"

type Switch interface {
	Run(ctx FlowContext) (r string, err error)
}

type SwitchFunc func(ctx FlowContext) (r string, err error)

func (this SwitchFunc) Run(ctx FlowContext) (r string, err error) { return this(ctx) }

type SwitchFlow struct {
	*CommonFlow
	flowSwitch Switch
	runnables  map[string]Runnable
}

func NewSwitchFlow() *SwitchFlow {
	o := &SwitchFlow{}
	o.CommonFlow = NewCommonFlow()
	o.runnables = make(map[string]Runnable)
	return o
}

func (this *SwitchFlow) AddPloy(cond string, ploy Ploy) *SwitchFlow {
	this.runnables[cond] = ploy
	this.CommonFlow.AddPloy(ploy)
	return this
}

func (this *SwitchFlow) AddPloyFunc(cond string, fn func(ctx FlowContext)) *SwitchFlow {
	this.AddPloy(cond, PloyFunc(fn))
	return this
}

func (this *SwitchFlow) AddFlow(cond string, flow Flow) *SwitchFlow {
	this.runnables[cond] = flow
	this.CommonFlow.AddFlow(flow)
	return this
}

func (this *SwitchFlow) AddHook(hooks ...Hook) *SwitchFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *SwitchFlow) SetSwitch(flowSwitch Switch) (r *SwitchFlow) {
	this.flowSwitch = flowSwitch
	return this
}
func (this *SwitchFlow) SetSwitchFunc(fn func(ctx FlowContext) (r string, err error)) (r *SwitchFlow) {
	this.SetSwitch(SwitchFunc(fn))
	return this
}

func (this *SwitchFlow) Run(ctx FlowContext) {
	if ctx.HasError() {
		return
	}

	if this.flowSwitch == nil {
		ctx.AddError(fmt.Errorf("switch flow must set switch first"))
		return
	}

	cond, err := this.flowSwitch.Run(ctx)
	if err != nil {
		ctx.AddError(fmt.Errorf("switch cond call error: %s", err))
		return
	}

	runnable, exists := this.runnables[cond]
	if !exists {
		ctx.AddError(fmt.Errorf("no such flow: %s", cond))
		return
	}

	this.hookMgr.before(ctx, runnable)
	runnable.Run(ctx)
	this.hookMgr.after(ctx, runnable)

	return
}
