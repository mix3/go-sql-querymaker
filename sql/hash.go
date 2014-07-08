package querymaker

type IHash interface {
	Keys() []string
	Value(string) interface{}
}

type Hash map[string]interface{}

func (h Hash) Keys() []string {
	keys := []string{}
	for k, _ := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h Hash) Value(key string) interface{} {
	if v, ok := h[key]; ok {
		return v
	}
	return nil
}

type OrderedHash struct {
	Hash  Hash
	Order []string
}

func NewOrderedHash(hash Hash, order ...string) IHash {
	return OrderedHash{
		Hash:  hash,
		Order: order,
	}
}

func (oh OrderedHash) Keys() []string {
	if oh.comp(oh.Hash.Keys(), oh.Order) {
		return oh.Order
	} else {
		return oh.Hash.Keys()
	}
}

func (oh OrderedHash) Value(key string) interface{} {
	if v, ok := oh.Hash[key]; ok {
		return v
	}
	return nil
}

func (h OrderedHash) comp(keys, order []string) bool {
	if len(keys) != len(order) {
		return false
	}

	m := make(map[string]int)
	for _, v := range append(keys, order...) {
		if _, ok := m[v]; ok {
			m[v]++
		} else {
			m[v] = 1
		}
	}

	for _, v := range m {
		if v < 2 {
			return false
		}
	}

	return true
}
