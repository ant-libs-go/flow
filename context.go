/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 16:20:17
# File Name: context.go
# Description:
####################################################################### */

package flow

import "sync"

type FlowContext interface {
	Lock()
	UnLock()
	RLock()
	RUnLock()
	AddError(err error)
	HasError() (r bool)
	GetError() (r error)
	GetErrors() (r []error)
}

type CoreContext struct {
	errors []error
	lock   sync.RWMutex
}

func (this *CoreContext) Lock()    { this.lock.Lock() }
func (this *CoreContext) UnLock()  { this.lock.Unlock() }
func (this *CoreContext) RLock()   { this.lock.RLock() }
func (this *CoreContext) RUnLock() { this.lock.RUnlock() }

func (this *CoreContext) AddError(err error) {
	this.Lock()
	defer this.UnLock()

	if err != nil {
		this.errors = append(this.errors, err)
	}
}

func (this *CoreContext) HasError() (r bool) {
	this.RLock()
	defer this.RUnLock()

	if len(this.errors) != 0 {
		r = true
	}
	return
}

func (this *CoreContext) GetError() (r error) {
	this.RLock()
	defer this.RUnLock()

	if this.HasError() {
		r = this.errors[0]
	}
	return
}

func (this *CoreContext) GetErrors() (r []error) {
	this.RLock()
	defer this.RUnLock()

	return this.errors
}
