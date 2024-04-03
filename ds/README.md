
- [Generic Data Structures](#generic-data-structures)
  - [Ring Buffer, Double Ended Queue](#ring-buffer-double-ended-queue)
  - [Queue](#queue)
  - [Stack](#stack)
  - [Heap](#heap)
  - [Set](#set)
  - [Benchmark](#benchmark)
- [Generic Iterator with Lazy Evaluation](#generic-iterator-with-lazy-evaluation)
  - [From X to Iterator](#from-x-to-iterator)
  - [Iterator Transformation](#iterator-transformation)
  - [Pipelines of Iterators](#pipelines-of-iterators)
  - [Drain Iterator](#drain-iterator)
  - [From Iterator to Y](#from-iterator-to-y)

## Generic Data Structures

### Ring Buffer, Double Ended Queue

[Ringbuf](ringbuf/ringbuf.go) is auto-scalable ring buffer. Ringbuf also implements deque (Double Ended Queue).

It's better to implement deque based on ringbuf instead of slice or list.

- In contrast to slice, ringbuf can reuse memory block of removed elements, and therefore less memory allocation and less GC.
    Another disadvantage of slice is that memory block can not be recycled in time. Even though a part of a slice's underlying array won't be reused, that part can not be recycled instantly, because the underlying array is referenced by the slice as a whole.
- In contrast to list, ringbuf has continuous memory space, no front and back pointers, and therefore higher access speed and less memory footprint.

**Import**

```go
import "github.com/burningxflame/gx/ds/ringbuf"
```

**API**:
PushBack, PushFront, PeekBack, PeekFront, PopBack, PopFront, Len, Range

**Sample**:
[ringbuf_test](ringbuf/ringbuf_test.go)

### Queue

[Queue](queue/queue.go): FIFO, a subset of deque.

**Import**

```go
import "github.com/burningxflame/gx/ds/queue"
```

**API**:
Enq, Peek, Deq, Len, Range

**Sample**:
[queue_test](queue/queue_test.go)

### Stack

[Stack](stack/stack.go): LIFO, another subset of deque.

**Import**

```go
import "github.com/burningxflame/gx/ds/stack"
```

**API**:
Push, Peek, Pop, Len, Range

**Sample**:
[stack_test](stack/stack_test.go)

### Heap

[Heap](heap/heap.go)

**Import**

```go
import "github.com/burningxflame/gx/ds/heap"
```

**API**:
Push, Pop, Len, Range

**Sample**:
[heap_test](heap/heap_test.go)

### Set

[Set](set/set.go)

**Import**

```go
import "github.com/burningxflame/gx/ds/set"
```

**API**:
Add, Delete, Len, Contain, Equal, Range, Union, Intersect, Diff

**Sample**:
[set_test](set/set_test.go)

### Benchmark

goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

Note: In each of the following benchmark results, the second colume (N) represents collection size, i.e., the number of elements in a collection.

**Ring Buffer, Double Ended Queue**:

```txt
BenchmarkPushBack-12      92285098         17.93 ns/op       47 B/op        0 allocs/op
BenchmarkPushFront-12     95364520         14.98 ns/op       46 B/op        0 allocs/op
BenchmarkPopFront-12      385959025          3.214 ns/op        0 B/op        0 allocs/op
BenchmarkPopBack-12       359027050          3.602 ns/op        0 B/op        0 allocs/op
```

**Queue**:

```txt
BenchmarkEnq-12    	58082024	        17.98 ns/op	      49 B/op	       0 allocs/op
BenchmarkDeq-12    	616338339	         4.427 ns/op	       0 B/op	       0 allocs/op
```

**Stack**:

```txt
BenchmarkPush-12    	86427447	        16.37 ns/op	      41 B/op	       0 allocs/op
BenchmarkPop-12     	585195102	         3.419 ns/op	       0 B/op	       0 allocs/op
```

**Heap**:

```txt
BenchmarkPush-12    	35654341	        41.86 ns/op	      50 B/op	       0 allocs/op
BenchmarkPop-12     	 5362929	       234.2 ns/op	       7 B/op	       0 allocs/op
```

**Set**:

```txt
BenchmarkAdd-12       	 7962375	       203.7 ns/op	      48 B/op	       0 allocs/op
BenchmarkDelete-12    	16687696	        90.03 ns/op	       0 B/op	       0 allocs/op
```

## Generic Iterator with Lazy Evaluation

In contrast to eager evaluation, where evaluation is performed immediately. An Iterator is lazily evaluated, i.e. not evaluated until it's necessary, e.g. when you drain an Iterator.

```go
import "github.com/burningxflame/gx/ds/iter"
```

### From X to Iterator

**Create Iterators from Builtin Types**

```go
// Create an Iterator from a slice
it := iter.FromSlice(l)

// Create an Iterator from a map
it := iter.FromMap(m)

// Create an Iterator from a channel
it := iter.FromChan(c)
```

**Create Iterators from Iters**
Iter is an interface representing a lazy iteration.

```go
// Represent a lazy iteration
type Iter[E any] interface {
	// Advance the iteration and return an Option wrapping the next value if exist.
	Next() Option[E]
}

// Represent an optional value
type Option[E any] struct {
	Val E    // the value if exist
	Ok  bool // true if exist, false otherwise
}
```

```go
// Create an Iterator consisting of all the elements of the first Iter, followed by all the elements of the second Iter, and so on.
it := iter.From(i1, i2, ...)
```

**Create Iterators from Rangers**
Ranger is an interface representing a normal iteration (not lazy).

```go
// Represent a normal iteration (not lazy)
type Ranger[E any] interface {
	// Iterate the collection and call fn for each element.
	ForEach(fn func(E))
}
```

```go
// Create an Iterator from a Ranger
it := iter.FromRanger(r)
```

**Samples**
[from_to](iter/from_to_test.go)

### Iterator Transformation

A transformation method transforms an Iterator in place, and returns the Iterator itself. This approach reduces memory footprint, especially for long pipelines of Iterators. And it enables method chaining as well, for easily building pipelines of Iterators.

```go
// Returns an Iterator consisting of those elements of the Iterator for which fn(e) returns true.
it2 := it.Filter(pred)

// Return an Iterator consisting of the results of applying fn to every element of the Iterator.
it2 := it.Map(fn)
```

[Go doesn't allow to define type parameters in methods](https://github.com/golang/go/issues/49085). So, to transform an Iterator of type E into an Iterator of another type F, use the function (not method) Map instead.

```go
func (it *Iterator[E]) Map(fn MapFn[E, E]) *Iterator[E]
vs
func Map[E, F any](it *Iterator[E], fn MapFn[E, F]) *Iterator[F]

// Map to another type
it2 := iter.Map(it, fn)
```

```go
// Return an Iterator consisting of the first n elements of the Iterator, or all elements if there are fewer than n.
it2 := it.Take(n)

// Returns an Iterator consisting of all but the first n elements of the Iterator.
it2 := it.Drop(n)

// Return an Iterator consisting of those elements of the Iterator as long as fn(e) returns true. Once fn(e) returns false, the rest of the elements are ignored.
it2 := it.TakeWhile(pred)

// Return an Iterator consisting of those elements of the Iterator starting from the first element for which fn(e) returns false.
it2 := it.DropWhile(pred)

// Return an Iterator consisting of all the elements of the first Iterator, followed by all the elements of the second Iterator, and so on.
it := iter.Chain(it1, it2, ...)
```

**Samples**
[transform](iter/transform_test.go)

### Pipelines of Iterators

It's easy to use method chaining to build pipelines of Iterators.

```go
it2 := it.Drop(n).Filter(pred).Map(fn).Take(m)...
```

Pipelines of Iterators are lazily evaluated, i.e. not evaluated until it's necessary, e.g. when you drain pipelines. This approach reduces memory footprint, especially for long pipelines of Iterators.

**Samples**
[pipeline](iter/transform_test.go)

### Drain Iterator

Draining an Iterator means completing iteration and producing a result or a side effect. This is the time when an Iterator is really evaluated. After draining, an Iterator is considered exhausted and it's pointless to use that Iterator again.

```go
// Call fn for each element of the Iterator.
it.ForEach(fn)

// Return the min element of the Iterator.
o := it.Min()

// Return the max element of the Iterator.
o := it.Max()

// Return a sorted slice of all element of the Iterator.
l := it.Sort(less)

// Return the result of applying fn to ini and the first element of the Iterator,
// then applying fn to that result and the second element, and so on.
// If the Iterator is empty, return ini and fn is not called.
acc := it.Reduce(ini, fn)
```

[Go doesn't allow to define type parameters in methods](https://github.com/golang/go/issues/49085). So, to reduce an Iterator of type E into a result of another type F, use the function (not method) Reduce instead.

```go
func (it *Iterator[E]) Reduce(ini E, fn func(acc E, e E) E) E
vs
func Reduce[A any, E any](it *Iterator[E], ini A, fn func(acc A, e E) A) A

// Reduce to another type
acc := iter.Reduce(it, ini, fn)
```

```go
// Call fn sequentially for each element of the Iterator. If fn returns false, stop iteration.
it.Range(fn)

// Return true if fn(e) is true for any element of the Iterator.
// If the Iterator is empty, return false.
ok := it.Any(pred)

// Return true if fn(e) is true for every element of the Iterator.
// If the Iterator is empty, return true.
ok := it.Every(pred)
```

**Samples**
[drain](iter/drain_test.go)

### From Iterator to Y

**Convert Iterators to Builtin Types**

```go
// Return a slice of all elements of the Iterator.
l := it.ToSlice()

// Return a map of all elements of the Iterator.
m := it.ToMap()

// Return a channel of all elements of the Iterator.
ch := it.ToChan()
```

**Convert Iterators to Collectors**
Collector is an interface representing a collector which collects all elements of an Iterator.

```go
// Represent a collector which collects all elements of an Iterator.
type Collector[E any] interface {
	// Add the element into the Collector
	Add(E)
}
```

```go
// Feed the Collector c with all elements of the Iterator.
it.To(c)
```

**Samples**
[from_to](iter/from_to_test.go)
