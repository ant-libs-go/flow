/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 15:02:25
# File Name: flow.go
# Description:
####################################################################### */

package flow

type Flow interface {
	Runnable
	setHookMgr(mgr *hookMgr)
}

type CommonFlow struct {
	runnables []Runnable
	hookMgr   *hookMgr
}

func NewCommonFlow() *CommonFlow {
	o := &CommonFlow{}
	o.runnables = make([]Runnable, 0, 10)
	return o
}

func (this *CommonFlow) AddPloy(ploy Ploy) {
	this.runnables = append(this.runnables, ploy)
}

func (this *CommonFlow) AddFlow(flow Flow) {
	flow.setHookMgr(this.hookMgr)
	this.runnables = append(this.runnables, flow)
}

func (this *CommonFlow) AddHook(hooks ...Hook) {
	if this.hookMgr == nil {
		this.setHookMgr(newHookMgr())
	}
	this.hookMgr.add(hooks)
}

func (this *CommonFlow) setHookMgr(mgr *hookMgr) {
	if this.hookMgr != nil {
		return
	}

	this.hookMgr = mgr
	for _, runnable := range this.runnables {
		if v, ok := runnable.(Flow); ok {
			v.setHookMgr(this.hookMgr)
		}
	}
}
