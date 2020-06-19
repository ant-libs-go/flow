/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 15:11:49
# File Name: flow.go
# Description:
####################################################################### */

package flow

func New(hooks ...Hook) *PipeFlow {
	o := NewPipeFlow()
	o.AddHook(hooks...)
	return o
}
