package filters

import . "reflect"

type Predicate			func(interface{}) bool

func Count(i interface{}, p interface{}) (n int) {
	switch i := i.(type) {
	case Enumerable:
		i.Each(p)
	case func() interface{}:
		p := p.(Predicate)
		for v := i(); v != nil; v = i() {
			if p(v) { n++ }
		}
	default:
		p := p.(Predicate)
		switch values := NewValue(i).(type) {
		case *ChanValue:
			for !values.Closed() {
				if p(values.Recv().Interface()) { n++ }
			}
		case *SliceValue:
			for j := 0; j < values.Len(); j++ {
				if p(values.Elem(j).Interface()) { n++ }
			}
		case *MapValue:
			for _, k := range values.Keys() {
				if p(values.Elem(k).Interface()) { n++ }
			}
		case *FuncValue:
			for {
				if v := values.Call([]Value{}); len(v) != 0 {
					p(v)
				} else {
					break
				}
			}
		default:
			panic("unsupported type")
		}
	}
	return
}


type Filterable interface {
	Count(p Predicate) int
	All(p Predicate) bool
	None(p Predicate) bool
}

type Container			[]interface{}
func (c Container) Repeat(count int) (n Container) {
	for ; count > 0; count-- {
		n = append(n, c...)
	}
	return
}
func (c Container) Count(p Predicate) (n int) {
	for _, v := range c {
		if p(v) {
			n++
		}
	}
	return
}
func (c Container) All(x Predicate) bool {
	for _, v := range c {
		if !x(v) {
			return false
		}
	}
	return true
}
func (c Container) None(x Predicate) bool {
	for _, v := range c {
		if x(v) {
			return false
		}
	}
	return true
}

func Any(c Container, p Predicate) bool {
	return !c.None(p)
}

type intPredicate		func(int) bool
func (i intPredicate) Prove(x interface{}) (status bool) {
	switch x := x.(type) {
	case int:
		status = i(x)
	case intContainer:
		for _, v := range x {
			if !i(v) {
				return false
			}
		}
		status = true
	}
	return
}

type intContainer		[]int
func (i intContainer) Prove(x intPredicate) bool {
	for _, v := range i {
		if !x(v) {
			return false
		}
	}
	return true
}
func (i intContainer) Repeat(count int) (n intContainer) {
	for ; count > 0; count-- {
		n = append(n, i...)
	}
	return
}

type floatPredicate		func(float) bool
type floatContainer		[]float
func (f floatContainer) Prove(x floatPredicate) bool {
	for _, v := range f {
		if !x(v) {
			return false
		}
	}
	return true
}

type uintPredicate		func(uint) bool
type uintContainer		[]uint
func (u uintContainer) Repeat(count int) (n uintContainer) {
	for ; count > 0; count-- {
		n = append(n, u...)
	}
	return
}
func (u uintContainer) Each(f interface{}) {
//	for _, v := range f {
		switch f := f.(type) {
		case func(uint):			for _, x := range u { f(x) }
		case func(interface{}):		for _, x := range u { f(x) }
		}
//	}
	return
}