package iter

import "reflect"

// Comparable values can compare against another value with the same type.
type Comparable interface {
	Less(Any) bool
}

func _less(x, y Any) bool {
	if c, ok := x.(Comparable); ok {
		return c.Less(y)
	}
	return reflectCompare(x, y) == -1
}

// Borrow from https://github.com/stretchr/testify/blob/master/assert/assertion_order.go
func reflectCompare(obj1, obj2 interface{}) int {
	k1 := reflect.ValueOf(obj1).Kind()
	k2 := reflect.ValueOf(obj2).Kind()
	if k2 != k1 {
		panic("cannot compare different types")
	}
	switch k1 {
	case reflect.Int:
		{
			intobj1 := obj1.(int)
			intobj2 := obj2.(int)
			if intobj1 > intobj2 {
				return -1
			}
			if intobj1 == intobj2 {
				return 0
			}
			if intobj1 < intobj2 {
				return 1
			}
		}
	case reflect.Int8:
		{
			int8obj1 := obj1.(int8)
			int8obj2 := obj2.(int8)
			if int8obj1 > int8obj2 {
				return -1
			}
			if int8obj1 == int8obj2 {
				return 0
			}
			if int8obj1 < int8obj2 {
				return 1
			}
		}
	case reflect.Int16:
		{
			int16obj1 := obj1.(int16)
			int16obj2 := obj2.(int16)
			if int16obj1 > int16obj2 {
				return -1
			}
			if int16obj1 == int16obj2 {
				return 0
			}
			if int16obj1 < int16obj2 {
				return 1
			}
		}
	case reflect.Int32:
		{
			int32obj1 := obj1.(int32)
			int32obj2 := obj2.(int32)
			if int32obj1 > int32obj2 {
				return -1
			}
			if int32obj1 == int32obj2 {
				return 0
			}
			if int32obj1 < int32obj2 {
				return 1
			}
		}
	case reflect.Int64:
		{
			int64obj1 := obj1.(int64)
			int64obj2 := obj2.(int64)
			if int64obj1 > int64obj2 {
				return -1
			}
			if int64obj1 == int64obj2 {
				return 0
			}
			if int64obj1 < int64obj2 {
				return 1
			}
		}
	case reflect.Uint:
		{
			uintobj1 := obj1.(uint)
			uintobj2 := obj2.(uint)
			if uintobj1 > uintobj2 {
				return -1
			}
			if uintobj1 == uintobj2 {
				return 0
			}
			if uintobj1 < uintobj2 {
				return 1
			}
		}
	case reflect.Uint8:
		{
			uint8obj1 := obj1.(uint8)
			uint8obj2 := obj2.(uint8)
			if uint8obj1 > uint8obj2 {
				return -1
			}
			if uint8obj1 == uint8obj2 {
				return 0
			}
			if uint8obj1 < uint8obj2 {
				return 1
			}
		}
	case reflect.Uint16:
		{
			uint16obj1 := obj1.(uint16)
			uint16obj2 := obj2.(uint16)
			if uint16obj1 > uint16obj2 {
				return -1
			}
			if uint16obj1 == uint16obj2 {
				return 0
			}
			if uint16obj1 < uint16obj2 {
				return 1
			}
		}
	case reflect.Uint32:
		{
			uint32obj1 := obj1.(uint32)
			uint32obj2 := obj2.(uint32)
			if uint32obj1 > uint32obj2 {
				return -1
			}
			if uint32obj1 == uint32obj2 {
				return 0
			}
			if uint32obj1 < uint32obj2 {
				return 1
			}
		}
	case reflect.Uint64:
		{
			uint64obj1 := obj1.(uint64)
			uint64obj2 := obj2.(uint64)
			if uint64obj1 > uint64obj2 {
				return -1
			}
			if uint64obj1 == uint64obj2 {
				return 0
			}
			if uint64obj1 < uint64obj2 {
				return 1
			}
		}
	case reflect.Float32:
		{
			float32obj1 := obj1.(float32)
			float32obj2 := obj2.(float32)
			if float32obj1 > float32obj2 {
				return -1
			}
			if float32obj1 == float32obj2 {
				return 0
			}
			if float32obj1 < float32obj2 {
				return 1
			}
		}
	case reflect.Float64:
		{
			float64obj1 := obj1.(float64)
			float64obj2 := obj2.(float64)
			if float64obj1 > float64obj2 {
				return -1
			}
			if float64obj1 == float64obj2 {
				return 0
			}
			if float64obj1 < float64obj2 {
				return 1
			}
		}
	case reflect.String:
		{
			stringobj1 := obj1.(string)
			stringobj2 := obj2.(string)
			if stringobj1 > stringobj2 {
				return -1
			}
			if stringobj1 == stringobj2 {
				return 0
			}
			if stringobj1 < stringobj2 {
				return 1
			}
		}
	}
	panic("unknown type")
}
