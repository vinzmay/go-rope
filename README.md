go-rope
=======

Description
-----------
[Go](http://www.golang.org) implementation of a persistent rope data structure, useful to store huge amounts of text to be operated upon. Persistent means that any operation on the rope doesn't modify it, so it's inherently thread safe.

Installation
------------
Package can be installed with go command:

	go get github.com/vinzmay/go-rope
	

Examples
--------

Creating a rope

```go
	//Empty rope
	r1 := new(rope.Rope)
	
	//Initialized rope
	r2 := rope.New("test rope")
```