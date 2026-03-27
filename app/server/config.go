package server

import (
	"go.uber.org/fx"
)

type ConfigProvider struct {
	// 把一个返回值，拆成多个依赖导出去
	// 返回值是 fx.Provide 里注册的构造函数执行出来的结果
	fx.Out
	Debug bool
}
