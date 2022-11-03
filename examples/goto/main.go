/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2022-11-03 19:23:47
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

type TestGoto struct {
	flow.Goto
}

func (this *TestGoto) Run(c flow.FlowContext) (r string, err error) {
	ctx := c.(*PContext)

	if ctx.Count == 10 {
		err = flow.BreakError
		return
	}

	if ctx.Count%2 == 0 {
		r = "test-a"
	} else {
		r = "test-b"
	}
	return
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New()

	gotoFlow := flow.NewGotoFlow().SetGoto(&TestGoto{})
	gotoFlow.AddPloy("test-a", &TestAPloy{})
	gotoFlow.AddPloy("test-b", &TestBPloy{})

	mainFlow.AddFlow(gotoFlow)

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
