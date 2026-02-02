package context

import (
	"fmt"
	"seshat/context"
	"slices"
)

type Queries struct {
	ctx *context.Ctx
	// queries are placed here
}

type Query[Key comparable, Val any] struct {
	// reference to context
	Context *context.Ctx

	// function description
	Func func(*context.Ctx, Key) Val
	Name string

	// caching per function
	Cache map[Key]Val

	// error checking
	Trace *[]context.Trace
}

func (q *Query[Key, Val]) Get(k Key) Val {
	if q == nil || q.Func == nil {
		panic(fmt.Sprintf("Query '%s' not defined!", q.Name))
	}

	return cacheQuery(q, k)
}

func cacheQuery[Key comparable, Val any](q *Query[Key, Val], k Key) Val {
	// caching
	if val, ok := q.Cache[k]; ok {
		return val
	}

	// cycle detection for Query system
	if q.Trace == nil {
		panic("Trace cannot be nil in Queries!")
	}

	if slices.ContainsFunc(*q.Trace, func(t context.Trace) bool { return t.Funcname == q.Name && t.Arg == k }) {
		panic(fmt.Sprintf("Caching did not work! Func '%s' with arg '%+v' was already in trace!", q.Name, k))
	}

	// actual function calling:
	res := q.Func(q.Context, k)

	*q.Trace = append(*q.Trace, context.Trace{Funcname: q.Name, Arg: k})

	q.Cache[k] = res
	return res
}
