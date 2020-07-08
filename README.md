# flow
代码解耦神器，支持并行、串行、分支解耦，增强代码复用性

# 功能
 - 将你的业务逻辑定义成独立的ploy，通过ParallelFlow、PipeFlow、SwitchFlow将业务逻辑进行串联
 - ParallelFlow: 并行逻辑
 - PipeFlow：顺序串行逻辑
 - SwitchFlow：分支逻辑

# 快速上手
 - 这里通过一个购物场景介绍flow该如何使用。由于场景是强行怼的，所以请不要在意场景是否合理以及异常处理。
 - 更多更详细的用法参考examples
```golang
// 常规写法
func Buy(userId int32, orderId int32) (r error) {
	user, _ := model.GetUserById(userId)
	order, _ := model.GetOrderById(orderId)

	if reason := userAntiSpam(user) ; reason != nil {
		return fmt.Errorf("user exception, %s", reason)
	}
	if reason := orderAntiSpam(order); reason != nil {
		return fmt.Errorf("order exception, %s", reason)
	}

	switch order.ProductType {
	case "CAR":
		r = Settle(user, BuyCar(user, order))
	case "COMPUTER":
		r = Settle(user, BuyComputer(user, order))
	case "SNACKS":
		r = Settle(user, BuySnacks(user, order))
	}
	return
}

// flow 的写法
type PContext struct {
	flow.CoreContext
	User *User
	Order *Order
}

// PipeFlow 的用法，顺序运行填充user和验证user的ploy
userFlow := flow.NewPipeFlow().AddPloy(new(FillingUser), new(UserAntiSpam))
orderFlow := flow.NewPipeFlow().AddPloy(new(FillingOrder), new(OrderAntiSpam))

// SwitchFlow 的用法，根据OrderProductTypeSwitch来判断后续运行的flow
buyFlow := flow.NewSwitchFlow().
	SetSwitch(new(OrderProductTypeSwitch)).
	// Settle Ploy 的复用
	AddFlow("CAR", flow.NewPipeFlow().AddPloy(new(BuyCar), new(Settle))).
	AddFlow("COMPUTER", flow.NewPipeFlow().AddPloy(new(BuyCar), new(Settle))).
	AddFlow("SNACKS", flow.NewPipeFlow().AddPloy(new(BuySnacks), new(Settle)))

// 主flow，user和order并行获取和验证
// 根据产品类型运行不同的购买逻辑以及结算逻辑
mainFlow := flow.New().
	AddFlow(flow.NewParallelFlow().AddFlow(userFlow, orderFlow)).
	AddFlow(buyFlow)

ctx := &PContext{}
mainFlow.Run(ctx)
```
