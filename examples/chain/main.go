/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-14 13:49:50
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

type TestSwitch struct {
	flow.Switch
}

func (this *TestSwitch) Run(c flow.FlowContext) (r string, err error) {
	//ctx := c.(*PContext)

	return "test-b", nil
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New().AddPloy(&TestAPloy{}, &TestBPloy{}).
		AddFlow(flow.NewPipeFlow().AddPloy(&TestAPloy{}, &TestBPloy{})).
		AddFlow(flow.NewParallelFlow().SetMaxProcess(2).AddPloy(&TestAPloy{}, &TestBPloy{})).
		AddFlow(flow.NewSwitchFlow().SetSwitch(&TestSwitch{}).AddPloy("test-a", &TestAPloy{}).AddPloy("test-b", &TestBPloy{}))

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
}
