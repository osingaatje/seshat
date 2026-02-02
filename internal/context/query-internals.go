package context

import (
	"fmt"
	"slices"
)

type Query[Key comparable, Val any] struct {
	// reference to context
	Context *Ctx

	// function description
	Func func(*Ctx, Key) Val
	Name string

	// caching per function
	Cache map[Key]Val

	// error checking
	Trace *[]Trace
}

func DefineQuery[K comparable, V any](c *Ctx, name string, f func(*Ctx, K) V) *Query[K, V] {
	fmt.Printf("Defining query '%s'", name)
	return &Query[K, V]{
		Context: c,
		Func:    f,
		Name:    name,
		Trace:   &c.Trace,
		Cache:   map[K]V{},
	}
}

func (q *Query[Key, Val]) Get(nameOfQuery string, k Key) Val {

	if q == nil || q.Func == nil {
		panic(fmt.Sprintf("Query '%s' not defined!", nameOfQuery))
	}

	fmt.Println("Get2")

	return cacheQuery(q, k)
}

func cacheQuery[Key comparable, Val any](q *Query[Key, Val], k Key) Val {
	fmt.Println("cache1")
	// caching
	if val, ok := q.Cache[k]; ok {
		return val
	}

	fmt.Println("cach2")
	// cycle detection for Query system
	if q.Trace == nil {
		panic("Trace cannot be nil in Queries!")
	}

	if slices.ContainsFunc(*q.Trace, func(t Trace) bool { return t.Funcname == q.Name && t.Arg == k }) {
		panic(fmt.Sprintf("Caching did not work! Func '%s' with arg '%+v' was already in trace!", q.Name, k))
	}

	// actual function calling:
	res := q.Func(q.Context, k)

	*q.Trace = append(*q.Trace, Trace{Funcname: q.Name, Arg: k})

	q.Cache[k] = res
	return res
}
