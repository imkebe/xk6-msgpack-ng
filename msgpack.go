// msgpack.go
package msgpack

import (
	"fmt"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"go.k6.io/k6/js/modules"
)

type (
	RootModule struct{}

	ModuleInstance struct {
		vu modules.VU
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

func init() {
	modules.Register("k6/x/msgpack", new(RootModule))
}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
	}
}

func (mi *ModuleInstance) Exports() modules.Exports {
	rt := mi.vu.Runtime()

	// Create the module object
	obj := rt.NewObject()

	// Define the serialize function
	_ = obj.Set("serialize", rt.ToValue(func(value interface{}) interface{} {
		// Marshal the value to MessagePack format
		data, err := msgpack.Marshal(value)
		if err != nil {
			panic(rt.NewGoError(fmt.Errorf("serialize error: %v", err)))
		}

		// Return the byte array as an ArrayBuffer
		return rt.NewArrayBuffer(data)
	}))

	// Define the deserialize function
	_ = obj.Set("deserialize", rt.ToValue(func(input interface{}) interface{} {
		var data []byte

		// Try to extract bytes from the input based on its type
		v := reflect.ValueOf(input)

		// Handle direct byte slice
		if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
			data = input.([]byte)
		} else {
			// Try to get the ArrayBuffer bytes using reflection
			// This is a fallback approach that tries to extract bytes from whatever object we get
			exported := rt.ToValue(input).Export()

			// Check if it's an object with a Bytes method
			if obj, ok := exported.(interface{ Bytes() []byte }); ok {
				data = obj.Bytes()
			} else {
				// Last resort: try to convert to string and use that
				data = []byte(fmt.Sprint(exported))
			}
		}

		// Unmarshal the data
		var result interface{}
		if err := msgpack.Unmarshal(data, &result); err != nil {
			panic(rt.NewGoError(fmt.Errorf("deserialize error: %v", err)))
		}

		return result
	}))

	return modules.Exports{
		Default: obj,
	}
}
