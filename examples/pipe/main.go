/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-13 19:55:13
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

	ctx.Count++

	time.Sleep(time.Second)
	fmt.Printf("%s: test-a\n", time.Now().Format("2006-01-02 15:04:05"))
}

type TestBPloy struct {
	flow.Ploy
}

func (TestBPloy) Run(c flow.FlowContext) {
	ctx := c.(*PContext)

	ctx.Count++

	time.Sleep(time.Second)
	fmt.Printf("%s: test-b\n", time.Now().Format("2006-01-02 15:04:05"))
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New()
	mainFlow.AddPloy(&TestAPloy{})
	mainFlow.AddPloy(&TestBPloy{})

	subFlow := flow.NewPipeFlow()
	subFlow.AddPloy(&TestAPloy{})
	subFlow.AddPloy(&TestBPloy{})
	mainFlow.AddFlow(subFlow)

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
}
