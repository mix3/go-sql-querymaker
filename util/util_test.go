package util

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if !reflect.DeepEqual(a, b) {
				t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
			}
		}
	}()
	if a != b {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestPop(t *testing.T) {
	arr := []interface{}{"a", "b", "c"}
	func() {
		a, b := Pop(arr)
		expect(t, a, "c")
		expect(t, b, []interface{}{"a", "b"})
		arr = b
	}()
	func() {
		a, b := Pop(arr)
		expect(t, a, "b")
		expect(t, b, []interface{}{"a"})
		arr = b
	}()
	func() {
		a, b := Pop(arr)
		expect(t, a, "a")
		expect(t, b, []interface{}{})
		arr = b
	}()
	func() {
		a, b := Pop(arr)
		expect(t, a, nil)
		expect(t, b, []interface{}{})
		arr = b
	}()
	func() {
		a, b := Pop(arr)
		expect(t, a, nil)
		expect(t, b, []interface{}{})
		arr = b
	}()
}

func TestPush(t *testing.T) {
	arr := []interface{}{}
	func() {
		act := Push(arr, "a")
		expect(t, act, []interface{}{"a"})
		arr = act
	}()
	func() {
		act := Push(arr, "b", "c")
		expect(t, act, []interface{}{"a", "b", "c"})
		arr = act
	}()
}

func TestShit(t *testing.T) {
	arr := []interface{}{"a", "b", "c"}
	func() {
		a, b := Shift(arr)
		expect(t, a, "a")
		expect(t, b, []interface{}{"b", "c"})
		arr = b
	}()
	func() {
		a, b := Shift(arr)
		expect(t, a, "b")
		expect(t, b, []interface{}{"c"})
		arr = b
	}()
	func() {
		a, b := Shift(arr)
		expect(t, a, "c")
		expect(t, b, []interface{}{})
		arr = b
	}()
	func() {
		a, b := Shift(arr)
		expect(t, a, nil)
		expect(t, b, []interface{}{})
		arr = b
	}()
	func() {
		a, b := Shift(arr)
		expect(t, a, nil)
		expect(t, b, []interface{}{})
		arr = b
	}()
}

func TestUnshift(t *testing.T) {
	arr := []interface{}{}
	func() {
		act := Unshift(arr, "a")
		expect(t, act, []interface{}{"a"})
		arr = act
	}()
	func() {
		act := Unshift(arr, "c", "b")
		expect(t, act, []interface{}{"c", "b", "a"})
		arr = act
	}()
}
