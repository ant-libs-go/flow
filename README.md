# Flow

代码解耦神器，支持并行、串行、分支解耦，增强代码复用性

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/flow?status.png)](http://godoc.org/github.com/ant-libs-go/flow)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/flow)](https://goreportcard.com/report/github.com/ant-libs-go/flow)

# 特性

* 将你的业务逻辑定义成独立的ploy，通过ParallelFlow、PipeFlow、SwitchFlow将业务逻辑进行串联
* ParallelFlow: 并行逻辑
* PipeFlow：顺序串行逻辑
* SwitchFlow：分支逻辑
* GotoFlow：Goto逻辑，根据条件在多个 ploy 中进行任意跳转，直到返回 BreakError
* GoroutineFlow：Goroutine并行逻辑，根据指定并发数运行指定的 ploy

## 安装

	go get github.com/ant-libs-go/flow

# 快速开始

* 这里通过一个购物场景介绍flow该如何使用。由于场景是强行怼的，所以请不要在意场景是否合理以及异常处理。
* 更多更详细的用法参考examples

```golang
// 常规写法
func Buy(userId int32, productId int32) (r error) {
	user, _ := model.GetUserById(userId)
	product, _ := model.GetProductById(productId)

	if reason := userAntiSpam(user) ; reason != nil {
		return fmt.Errorf("user exception, %s", reason)
	}
	if reason := productAntiSpam(product); reason != nil {
		return fmt.Errorf("product exception, %s", reason)
	}

	switch product.Type {
	case "CAR":
		r = Settle(user, BuyCar(user, product))
	case "COMPUTER":
		r = Settle(user, BuyComputer(user, product))
	case "SNACKS":
		r = Settle(user, BuySnacks(user, product))
	}
	return
}

// flow 的写法
type PContext struct {
	flow.CoreContext
	User *User
	Product *Product
}

// PipeFlow 的用法，顺序运行填充user和验证user的ploy
userFlow := flow.NewPipeFlow().AddPloy(new(FillingUser), new(UserAntiSpam))
productFlow := flow.NewPipeFlow().AddPloy(new(FillingProduct), new(ProductAntiSpam))

// SwitchFlow 的用法，根据ProductTypeSwitch来判断后续运行的flow
buyFlow := flow.NewSwitchFlow().
	SetSwitch(new(ProductTypeSwitch)).
	// Settle Ploy 的复用
	AddFlow("CAR", flow.NewPipeFlow().AddPloy(new(BuyCar), new(Settle))).
	AddFlow("COMPUTER", flow.NewPipeFlow().AddPloy(new(BuyComputer), new(Settle))).
	AddFlow("SNACKS", flow.NewPipeFlow().AddPloy(new(BuySnacks), new(Settle)))

// 主flow，user和product并行获取和验证
// 根据产品类型运行不同的购买逻辑以及结算逻辑
mainFlow := flow.New().
	AddFlow(flow.NewParallelFlow().AddFlow(userFlow, productFlow)).
	AddFlow(buyFlow)

ctx := &PContext{}
mainFlow.Run(ctx)
```

# 高级用法

* 使用 https://github.com/ant-libs-go/pool 对`PContext`进行复用

