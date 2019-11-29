package iter

import "reflect"

// Any represents any type.
type Any interface{}

func _eq(x, y Any) bool {
	type ieq interface{ Eq(Any) bool }
	if e, ok := x.(ieq); ok {
		return e.Eq(y)
	}
	return x == y
}

func _ne(x, y Any) bool {
	return !_eq(x, y)
}

func _less(x, y Any) bool {
	type iless interface{ Less(Any) bool }
	if c, ok := x.(iless); ok {
		return c.Less(y)
	}
	return _cmp(x, y) == -1
}

func _cmp(x, y Any) int {
	type icmp interface{ Cmp(Any) int }
	if t, ok := x.(icmp); ok {
		return t.Cmp(y)
	}
	return reflectCompare(x, y)
}

func _inc(x Any) Any {
	type iinc interface{ Inc() Any }
	if i, ok := x.(iinc); ok {
		return i.Inc()
	}
	return reflectInc(x)
}

// Borrow from https://github.com/stretchr/testify/blob/master/assert/assertion_order.go
func reflectCompare(obj1, obj2 Any) int {
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
				return 1
			}
			if intobj1 == intobj2 {
				return 0
			}
			if intobj1 < intobj2 {
				return -1
			}
		}
	case reflect.Int8:
		{
			int8obj1 := obj1.(int8)
			int8obj2 := obj2.(int8)
			if int8obj1 > int8obj2 {
				return 1
			}
			if int8obj1 == int8obj2 {
				return 0
			}
			if int8obj1 < int8obj2 {
				return -1
			}
		}
	case reflect.Int16:
		{
			int16obj1 := obj1.(int16)
			int16obj2 := obj2.(int16)
			if int16obj1 > int16obj2 {
				return 1
			}
			if int16obj1 == int16obj2 {
				return 0
			}
			if int16obj1 < int16obj2 {
				return -1
			}
		}
	case reflect.Int32:
		{
			int32obj1 := obj1.(int32)
			int32obj2 := obj2.(int32)
			if int32obj1 > int32obj2 {
				return 1
			}
			if int32obj1 == int32obj2 {
				return 0
			}
			if int32obj1 < int32obj2 {
				return -1
			}
		}
	case reflect.Int64:
		{
			int64obj1 := obj1.(int64)
			int64obj2 := obj2.(int64)
			if int64obj1 > int64obj2 {
				return 1
			}
			if int64obj1 == int64obj2 {
				return 0
			}
			if int64obj1 < int64obj2 {
				return -1
			}
		}
	case reflect.Uint:
		{
			uintobj1 := obj1.(uint)
			uintobj2 := obj2.(uint)
			if uintobj1 > uintobj2 {
				return 1
			}
			if uintobj1 == uintobj2 {
				return 0
			}
			if uintobj1 < uintobj2 {
				return -1
			}
		}
	case reflect.Uint8:
		{
			uint8obj1 := obj1.(uint8)
			uint8obj2 := obj2.(uint8)
			if uint8obj1 > uint8obj2 {
				return 1
			}
			if uint8obj1 == uint8obj2 {
				return 0
			}
			if uint8obj1 < uint8obj2 {
				return -1
			}
		}
	case reflect.Uint16:
		{
			uint16obj1 := obj1.(uint16)
			uint16obj2 := obj2.(uint16)
			if uint16obj1 > uint16obj2 {
				return 1
			}
			if uint16obj1 == uint16obj2 {
				return 0
			}
			if uint16obj1 < uint16obj2 {
				return -1
			}
		}
	case reflect.Uint32:
		{
			uint32obj1 := obj1.(uint32)
			uint32obj2 := obj2.(uint32)
			if uint32obj1 > uint32obj2 {
				return 1
			}
			if uint32obj1 == uint32obj2 {
				return 0
			}
			if uint32obj1 < uint32obj2 {
				return -1
			}
		}
	case reflect.Uint64:
		{
			uint64obj1 := obj1.(uint64)
			uint64obj2 := obj2.(uint64)
			if uint64obj1 > uint64obj2 {
				return 1
			}
			if uint64obj1 == uint64obj2 {
				return 0
			}
			if uint64obj1 < uint64obj2 {
				return -1
			}
		}
	case reflect.Float32:
		{
			float32obj1 := obj1.(float32)
			float32obj2 := obj2.(float32)
			if float32obj1 > float32obj2 {
				return 1
			}
			if float32obj1 == float32obj2 {
				return 0
			}
			if float32obj1 < float32obj2 {
				return -1
			}
		}
	case reflect.Float64:
		{
			float64obj1 := obj1.(float64)
			float64obj2 := obj2.(float64)
			if float64obj1 > float64obj2 {
				return 1
			}
			if float64obj1 == float64obj2 {
				return 0
			}
			if float64obj1 < float64obj2 {
				return -1
			}
		}
	case reflect.String:
		{
			stringobj1 := obj1.(string)
			stringobj2 := obj2.(string)
			if stringobj1 > stringobj2 {
				return 1
			}
			if stringobj1 == stringobj2 {
				return 0
			}
			if stringobj1 < stringobj2 {
				return -1
			}
		}
	}
	panic("unknown type")
}

func reflectInc(x Any) Any {
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
	}
	panic("unknown type")
}
