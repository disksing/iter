package iter

import (
	"fmt"
	"reflect"
	"strings"
)

func reflectCmp(obj1, obj2 any) int {
	k1 := reflect.ValueOf(obj1).Kind()
	k2 := reflect.ValueOf(obj2).Kind()
	if k2 != k1 {
		panic("cannot sub different types")
	}
	switch k1 {
	case reflect.Int:
		switch {
		case obj1.(int) == obj2.(int):
			return 0
		case obj1.(int) < obj2.(int):
			return -1
		default:
			return 1
		}
	case reflect.Int8:
		switch {
		case obj1.(int8) == obj2.(int8):
			return 0
		case obj1.(int8) < obj2.(int8):
			return -1
		default:
			return 1
		}
	case reflect.Int16:
		switch {
		case obj1.(int16) == obj2.(int16):
			return 0
		case obj1.(int16) < obj2.(int16):
			return -1
		default:
			return 1
		}
	case reflect.Int32:
		switch {
		case obj1.(int32) == obj2.(int32):
			return 0
		case obj1.(int32) < obj2.(int32):
			return -1
		default:
			return 1
		}
	case reflect.Int64:
		switch {
		case obj1.(int64) == obj2.(int64):
			return 0
		case obj1.(int64) < obj2.(int64):
			return -1
		default:
			return 1
		}
	case reflect.Uint:
		switch {
		case obj1.(uint) == obj2.(uint):
			return 0
		case obj1.(uint) < obj2.(uint):
			return -1
		default:
			return 1
		}
	case reflect.Uint8:
		switch {
		case obj1.(uint8) == obj2.(uint8):
			return 0
		case obj1.(uint8) < obj2.(uint8):
			return -1
		default:
			return 1
		}
	case reflect.Uint16:
		switch {
		case obj1.(uint16) == obj2.(uint16):
			return 0
		case obj1.(uint16) < obj2.(uint16):
			return -1
		default:
			return 1
		}
	case reflect.Uint32:
		switch {
		case obj1.(uint32) == obj2.(uint32):
			return 0
		case obj1.(uint32) < obj2.(uint32):
			return -1
		default:
			return 1
		}
	case reflect.Uint64:
		switch {
		case obj1.(uint64) == obj2.(uint64):
			return 0
		case obj1.(uint64) < obj2.(uint64):
			return -1
		default:
			return 1
		}
	case reflect.Float32:
		switch {
		case obj1.(float32) == obj2.(float32):
			return 0
		case obj1.(float32) < obj2.(float32):
			return -1
		default:
			return 1
		}
	case reflect.Float64:
		switch {
		case obj1.(float64) == obj2.(float64):
			return 0
		case obj1.(float64) < obj2.(float64):
			return -1
		default:
			return 1
		}
	case reflect.String:
		return strings.Compare(obj1.(string), obj2.(string))
	}
	panic("unknown type")
}

func reflectInc(x any) any {
	switch v := x.(type) {
	case int:
		return v + 1
	case int8:
		return v + 1
	case int16:
		return v + 1
	case int32:
		return v + 1
	case int64:
		return v + 1
	case uint:
		return v + 1
	case uint8:
		return v + 1
	case uint16:
		return v + 1
	case uint32:
		return v + 1
	case uint64:
		return v + 1
	case float32:
		return v + 1
	case float64:
		return v + 1
	case complex64:
		return v + 1
	case complex128:
		return v + 1
	}
	panic("unknown type")
}

func reflectAdd(obj1, obj2 any) any {
	k1 := reflect.ValueOf(obj1).Kind()
	k2 := reflect.ValueOf(obj2).Kind()
	if k2 != k1 {
		panic(fmt.Sprintf("cannot add different types: %v vs %v", k1, k2))
	}
	switch k1 {
	case reflect.Int:
		return obj1.(int) + obj2.(int)
	case reflect.Int8:
		return obj1.(int8) + obj2.(int8)
	case reflect.Int16:
		return obj1.(int16) + obj2.(int16)
	case reflect.Int32:
		return obj1.(int32) + obj2.(int32)
	case reflect.Int64:
		return obj1.(int64) + obj2.(int64)
	case reflect.Uint:
		return obj1.(uint) + obj2.(uint)
	case reflect.Uint8:
		return obj1.(uint8) + obj2.(uint8)
	case reflect.Uint16:
		return obj1.(uint16) + obj2.(uint16)
	case reflect.Uint32:
		return obj1.(uint32) + obj2.(uint32)
	case reflect.Uint64:
		return obj1.(uint64) + obj2.(uint64)
	case reflect.Float32:
		return obj1.(float32) + obj2.(float32)
	case reflect.Float64:
		return obj1.(float64) + obj2.(float64)
	case reflect.Complex64:
		return obj1.(complex64) + obj2.(complex64)
	case reflect.Complex128:
		return obj1.(complex128) + obj2.(complex128)
	case reflect.String:
		return obj1.(string) + obj2.(string)
	}
	panic("unknown type")
}

func reflectSub(obj1, obj2 any) any {
	k1 := reflect.ValueOf(obj1).Kind()
	k2 := reflect.ValueOf(obj2).Kind()
	if k2 != k1 {
		panic("cannot sub different types")
	}
	switch k1 {
	case reflect.Int:
		return obj1.(int) - obj2.(int)
	case reflect.Int8:
		return obj1.(int8) - obj2.(int8)
	case reflect.Int16:
		return obj1.(int16) - obj2.(int16)
	case reflect.Int32:
		return obj1.(int32) - obj2.(int32)
	case reflect.Int64:
		return obj1.(int64) - obj2.(int64)
	case reflect.Uint:
		return obj1.(uint) - obj2.(uint)
	case reflect.Uint8:
		return obj1.(uint8) - obj2.(uint8)
	case reflect.Uint16:
		return obj1.(uint16) - obj2.(uint16)
	case reflect.Uint32:
		return obj1.(uint32) - obj2.(uint32)
	case reflect.Uint64:
		return obj1.(uint64) - obj2.(uint64)
	case reflect.Float32:
		return obj1.(float32) - obj2.(float32)
	case reflect.Float64:
		return obj1.(float64) - obj2.(float64)
	case reflect.Complex64:
		return obj1.(complex64) - obj2.(complex64)
	case reflect.Complex128:
		return obj1.(complex128) - obj2.(complex128)
	}
	panic("unknown type")
}

func reflectMul(obj1, obj2 any) any {
	k1 := reflect.ValueOf(obj1).Kind()
	k2 := reflect.ValueOf(obj2).Kind()
	if k2 != k1 {
		panic("cannot muliply different types")
	}
	switch k1 {
	case reflect.Int:
		return obj1.(int) * obj2.(int)
	case reflect.Int8:
		return obj1.(int8) * obj2.(int8)
	case reflect.Int16:
		return obj1.(int16) * obj2.(int16)
	case reflect.Int32:
		return obj1.(int32) * obj2.(int32)
	case reflect.Int64:
		return obj1.(int64) * obj2.(int64)
	case reflect.Uint:
		return obj1.(uint) * obj2.(uint)
	case reflect.Uint8:
		return obj1.(uint8) * obj2.(uint8)
	case reflect.Uint16:
		return obj1.(uint16) * obj2.(uint16)
	case reflect.Uint32:
		return obj1.(uint32) * obj2.(uint32)
	case reflect.Uint64:
		return obj1.(uint64) * obj2.(uint64)
	case reflect.Float32:
		return obj1.(float32) * obj2.(float32)
	case reflect.Float64:
		return obj1.(float64) * obj2.(float64)
	case reflect.Complex64:
		return obj1.(complex64) * obj2.(complex64)
	case reflect.Complex128:
		return obj1.(complex128) * obj2.(complex128)
	}
	panic("unknown type")
}
