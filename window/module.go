package window

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

const ModuleName = "window"

func init() {
	require.RegisterCoreModule(ModuleName, Require)
}

func Enable(runtime *goja.Runtime) {
	runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}

func Require(runtime *goja.Runtime, module *goja.Object) {
	func(runtime *goja.Runtime, module *goja.Object) {
		c := &fetchMode{
			r: runtime,
		}
		o := module.Get("exports").(*goja.Object)
		o.Set("fetch", c.Fetch)
	}(runtime, module)
}
