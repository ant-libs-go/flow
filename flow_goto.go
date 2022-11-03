/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-08-19 13:39:06
# File Name: flow_goto.go
# Description:
####################################################################### */

package flow

import (
	"errors"
	"fmt"
)

var BreakError = errors.New("break error")

type Goto interface {
	Run(ctx FlowContext) (r string, err error)
}

type GotoFunc func(ctx FlowContext) (r string, err error)

func (this GotoFunc) Run(ctx FlowContext) (r string, err error) { return this(ctx) }

type GotoFlow struct {
	*CommonFlow
	flowGoto  Goto
	runnables map[string][]Runnable
}

func NewGotoFlow() *GotoFlow {
	o := &GotoFlow{}
	o.CommonFlow = NewCommonFlow()
	o.runnables = make(map[string][]Runnable)
	return o
}

func (this *GotoFlow) AddPloy(cond string, ploys ...Ploy) *GotoFlow {
	for _, ploy := range ploys {
		this.runnables[cond] = append(this.runnables[cond], ploy)
		this.CommonFlow.AddPloy(ploy)
	}
	return this
}

func (this *GotoFlow) AddPloyFunc(cond string, fns ...func(ctx FlowContext)) *GotoFlow {
	for _, fn := range fns {
		this.AddPloy(cond, PloyFunc(fn))
	}
	return this
}

func (this *GotoFlow) AddFlow(cond string, flows ...Flow) *GotoFlow {
	for _, flow := range flows {
		this.runnables[cond] = append(this.runnables[cond], flow)
		this.CommonFlow.AddFlow(flow)
	}
	return this
}

func (this *GotoFlow) AddHook(hooks ...Hook) *GotoFlow {
	this.CommonFlow.AddHook(hooks...)
	return this
}

func (this *GotoFlow) SetGoto(flowGoto Goto) (r *GotoFlow) {
	this.flowGoto = flowGoto
	return this
}
func (this *GotoFlow) SetGotoFunc(fn func(ctx FlowContext) (r string, err error)) (r *GotoFlow) {
	this.SetGoto(GotoFunc(fn))
	return this
}

func (this *GotoFlow) Run(ctx FlowContext) {
	for {
		if this.flowGoto == nil {
			ctx.AddError(fmt.Errorf("goto flow must set goto first"))
			return
		}

		cond, err := this.flowGoto.Run(ctx)
		if err == BreakError {
			return
		}
		if err != nil {
			ctx.AddError(fmt.Errorf("goto cond call error: %s", err))
			return
		}

		runnables, ok := this.runnables[cond]
		if !ok {
			ctx.AddError(fmt.Errorf("no such flow: %s", cond))
			return
		}

		for _, runnable := range runnables {
			if ctx.HasError() {
				break
			}
			this.hookMgr.before(ctx, runnable)
			runnable.Run(ctx)
			this.hookMgr.after(ctx, runnable)
		}
	}
}
