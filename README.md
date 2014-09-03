go-rope v0.5
============

Description
-----------

[Go](http://www.golang.org) implementation of a persistent [rope](http://en.wikipedia.org/wiki/Rope_%28data_structure%29) data structure, useful to manipulate large text. Persistent means that any operation on the rope doesn't modify it, so it's inherently thread safe.

TODO: Rebalancing

Installation
------------

Package can be installed with go command:

	go get github.com/vinzmay/go-rope
	

Examples
--------

## Create a rope

```go
	//Empty rope
	r1 := new(rope.Rope)
	
	//Initialized rope
	r2 := rope.New("test rope")
```

## Concatenate two ropes

```go
	r1 := rope.New("abc")
	r2 := rope.New("def")
	
	r3 := r1.Concat(r2) // "abcdef"
```

## Split a rope

```go
	r1 := rope.New("abcdef")
	
	r2, r3 := r1.Split(4) // "abcd", "ef"
```

## Delete from a rope

```go
	r1 := rope.New("abcdef")
	
	r2 := r1.Delete(3, 2) // "abef"
```

## Insert in a rope

```go
	r1 := rope.New("abcdef")
	
	r2 := r1.Insert(3, "xxx") // "abcxxxdef"
```

## Rope substring

```go
	r1 := rope.New("abcdef")
	
	r2 := r1.Substr(3, 2) // "cd"
```
