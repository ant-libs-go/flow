/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 15:03:26
# File Name: runnable.go
# Description:
####################################################################### */

package flow

type Runnable interface {
	Run(ctx FlowContext)
}
