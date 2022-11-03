/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2022-11-03 19:28:19
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

type TestGoroutine struct {
	flow.Goroutine
}

func (this *TestGoroutine) Run(c flow.FlowContext) (r int, err error) {
	//ctx := c.(*PContext)

	return 5, nil
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New()

	goroutineFlow := flow.NewGoroutineFlow().SetGoroutine(&TestGoroutine{})
	goroutineFlow.AddPloy(&TestAPloy{})
	goroutineFlow.AddPloy(&TestBPloy{})

	mainFlow.AddFlow(goroutineFlow)

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
