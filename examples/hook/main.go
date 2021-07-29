/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-06-14 13:01:51
# File Name: main.go
# Description:
####################################################################### */

package main

import (
	"fmt"
	"reflect"
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

type TestCPloy struct {
	flow.Ploy
}

func (TestCPloy) Run(c flow.FlowContext) {
	ctx := c.(*PContext)

	ctx.Count++

	time.Sleep(time.Second)
	fmt.Printf("%s: test-c\n", time.Now().Format("2006-01-02 15:04:05"))
}

type TestSwitch struct {
	flow.Switch
}

func (this *TestSwitch) Run(c flow.FlowContext) (r string, err error) {
	//ctx := c.(*PContext)

	return "test-b", nil
}

type TestHook struct {
	flow.Hook
}

func (this *TestHook) Before(c flow.FlowContext, runnable flow.Runnable) {
	//ctx := c.(*PContext)

	fmt.Println("before", reflect.TypeOf(runnable))
}

func (this *TestHook) After(c flow.FlowContext, runnable flow.Runnable) {
	//ctx := c.(*PContext)

	fmt.Println("after", reflect.TypeOf(runnable))
}

type TestHook2 struct {
	flow.Hook
}

func (this *TestHook2) Before(c flow.FlowContext, runnable flow.Runnable) {
	//ctx := c.(*PContext)

	fmt.Println("before2", reflect.TypeOf(runnable))
}

func (this *TestHook2) After(c flow.FlowContext, runnable flow.Runnable) {
	//ctx := c.(*PContext)

	fmt.Println("after2", reflect.TypeOf(runnable))
}

type PContext struct {
	flow.CoreContext
	Count int
}

func main() {
	mainFlow := flow.New()
	mainFlow.AddHook(&TestHook{})
	mainFlow.AddPloy(&TestAPloy{}, &TestBPloy{})

	subFlow := flow.NewPipeFlow()
	subFlow.AddPloy(&TestAPloy{}, &TestBPloy{})

	parallelFlow := flow.NewParallelFlow().SetMaxProcess(2)
	parallelFlow.AddPloy(&TestAPloy{}, &TestCPloy{})

	switchFlow := flow.NewSwitchFlow().SetSwitch(&TestSwitch{})
	switchFlow.AddPloy("test-a", &TestAPloy{})
	switchFlow.AddPloy("test-b", &TestBPloy{})
	switchFlow.AddHook(&TestHook2{})

	mainFlow.AddFlow(subFlow, parallelFlow, switchFlow)

	ctx := &PContext{}
	mainFlow.Run(ctx)
	fmt.Println(ctx.Count)
	fmt.Println(ctx.GetErrors())
}
