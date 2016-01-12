package poloniex

import (
	"sort"
	"sync"

	"github.com/deckarep/golang-set"
)

type (
	OrderedFloatMap struct {
		m map[float64]float64
		k mapset.Set
		s sync.Mutex
	}
)

func NewOrderedFloatMap() (om *OrderedFloatMap) {
	om = &OrderedFloatMap{}
	om.m = map[float64]float64{}
	om.k = mapset.NewSet()
	om.s = sync.Mutex{}
	return
}

func (o *OrderedFloatMap) Set(k, v float64) {
	o.s.Lock()
	defer o.s.Unlock()
	o.m[k] = v + o.m[k]
	o.k.Add(k)
}

func (o *OrderedFloatMap) Get(k float64) float64 {
	return o.m[k]
}

func (o *OrderedFloatMap) Del(k float64) float64 {
	v := o.m[k]
	o.k.Remove(k)
	delete(o.m, k)
	return v
}

func (o *OrderedFloatMap) Inc(k, v float64) float64 {
	o.m[k] = o.m[k] + v
	return o.m[k]
}

func (o *OrderedFloatMap) Dec(k, v float64) float64 {
	o.m[k] = o.m[k] - v
	return o.m[k]
}

func (o *OrderedFloatMap) Contains(k float64) bool {
	return o.k.Contains(k)
}

func (o *OrderedFloatMap) Keys() []float64 {
	f := []float64{}
	for k := range o.k.Iter() {
		f = append(f, k.(float64))
	}
	sort.Float64s(f)
	return f
}

func (o *OrderedFloatMap) ReverseKeys() []float64 {
	f := []float64{}
	for k := range o.k.Iter() {
		f = append(f, k.(float64))
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(f)))
	return f
}
