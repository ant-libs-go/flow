/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-14 11:48:25
# File Name: main.go
# Description:
####################################################################### */

package main

import (
	"fmt"
	"time"

	"github.com/ant-libs-go/flow"
)

type TestAPloy struct {
	flow.Ploy
}

func (TestAPloy) Run(c flow.FlowContext) {
	ctx := c.(*PContext)

	ctx.Lock()
	ctx.Count++
	ctx.UnLock()

	time.Sleep(time.Second)
	fmt.Printf("%s: test-a\n", time.Now().Format("2006-01-02 15:04:05"))
}

type TestBPloy struct {
	flow.Ploy
}

func (TestBPloy) Run(c flow.FlowContext) {
	ctx := c.(*PContext)

	ctx.Lock()
	ctx.Count++
	ctx.UnLock()

	time.Sleep(time.Second)
	fmt.Printf("%s: test-b\n", time.Now().Format("2006-01-02 15:04:05"))
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New()

	parallelFlow := flow.NewParallelFlow().SetMaxProcess(2)
	parallelFlow.AddPloy(&TestAPloy{})
	parallelFlow.AddPloy(&TestBPloy{})

	mainFlow.AddFlow(parallelFlow)

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
}
