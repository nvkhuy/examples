package helper

import (
	"reflect"
)

// ToPtr wraps the given value with pointer: V => *V, *V => **V, etc.
func ToPtr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}
