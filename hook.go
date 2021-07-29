/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 14:59:19
# File Name: hook.go
# Description:
####################################################################### */

package flow

type Hook interface {
	Before(ctx FlowContext, runable Runnable)
	After(ctx FlowContext, runable Runnable)
}

type hookMgr struct {
	hooks []Hook
}

func newHookMgr() *hookMgr {
	o := &hookMgr{}
	o.hooks = make([]Hook, 0, 10)
	return o
}

func (this *hookMgr) add(hooks []Hook) {
	for _, hook := range hooks {
		this.hooks = append(this.hooks, hook)
	}
}

func (this *hookMgr) before(ctx FlowContext, runable Runnable) {
	if this == nil {
		return
	}
	for _, hook := range this.hooks {
		hook.Before(ctx, runable)
	}
}

func (this *hookMgr) after(ctx FlowContext, runable Runnable) {
	if this == nil {
		return
	}
	for _, hook := range this.hooks {
		hook.After(ctx, runable)
	}
}
