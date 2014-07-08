package util

func Pop(arr []interface{}) (interface{}, []interface{}) {
	switch t := len(arr); {
	case t == 0:
		return nil, []interface{}{}
	case t == 1:
		return arr[0], []interface{}{}
	default:
		return arr[t-1], arr[:t-1]
	}
}

func Push(arr []interface{}, in ...interface{}) []interface{} {
	return append(arr, in...)
}

func Shift(arr []interface{}) (interface{}, []interface{}) {
	switch t := len(arr); {
	case t == 0:
		return nil, []interface{}{}
	case t == 1:
		return arr[0], []interface{}{}
	default:
		return arr[0], arr[1:]
	}
}

func Unshift(arr []interface{}, in ...interface{}) []interface{} {
	return append(in, arr...)
}
