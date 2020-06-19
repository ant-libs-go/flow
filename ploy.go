/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 15:04:08
# File Name: ploy.go
# Description:
####################################################################### */

package flow

type Ploy Runnable

type PloyFunc func(ctx FlowContext)

func (this PloyFunc) Run(ctx FlowContext) { this(ctx) }
