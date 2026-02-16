package context

//----------------------------------------//
// Logic for defining and caching queries //
//----------------------------------------//

import (
	"fmt"
	"slices"
)

type Query[K comparable, Val any] struct {
	// reference to context
	Context *Ctx

	// function description
	Func func(*Ctx, K) Val
	Name string

	// caching per function
	Cache map[K]Val
}

func DefineQuery[K comparable, V any](c *Ctx, name string, f func(*Ctx, K) V) *Query[K, V] {
	//c.LogDebug("Defining query '%s'\n", name)
	return &Query[K, V]{
		Context: c,
		Func:    f,
		Name:    name,
		Cache:   map[K]V{},
	}
}

func (q *Query[K, V]) Get(nameOfQuery string, k K) V {

	if q == nil || q.Func == nil {
		panic(fmt.Sprintf("Query '%s' not defined!", nameOfQuery))
	}

	return cacheQuery(q, k)
}

func cacheQuery[K comparable, V any](q *Query[K, V], k K) V {
	if q == nil {
		panic("Query was nil in cacheQuery!")
	}

	// caching
	if val, ok := q.Cache[k]; ok {
		return val
	}

	ctx := q.Context
	var trace *[]Trace = &ctx.Trace

	// detect query loops
	if slices.ContainsFunc(*trace, func(t Trace) bool { return t.Funcname == q.Name && t.Arg == k }) {
		panic(fmt.Sprintf("Query loop! Caching did not work or your queries are ill-defined. Func '%s' with arg '%+v' was already in trace!", q.Name, k))
	}

	// actual function calling:
	res := q.Func(q.Context, k)

	// make sure we log that this query has been called and with which args.
	*trace = append(*trace, Trace{Funcname: q.Name, Arg: k})

	q.Cache[k] = res
	return res
}
