package main

import (
	"fmt"
	"github.com/vinzmay/go-rope"
)

func main() {
	a := rope.New("Hello, 世界")
	fmt.Println("a:", a)
	fmt.Println("a len:", a.Len())
	//fmt.Println("a json:", a.ToJSON())
	a1, a2 := a.Split(7)
	fmt.Println("a1:", a1)
	fmt.Println("a2:", a2)
	fmt.Println("a_test:", a1.Concat(a2))
	fmt.Println("a[8]:", string(a.Index(8)))
}

func main2() {
	e := rope.New("Hello_")
	f := rope.New("my_")
	c := e.Concat(f)
	j := rope.New("na")
	k := rope.New("me_i")
	g := j.Concat(k)
	m := rope.New("s")
	n := rope.New("_Simon")
	h := m.Concat(n)
	d := g.Concat(h)
	b := c.Concat(d)
	a := b.Concat(nil)
	a1, a2 := a.Split(7)
	a3 := a.Insert(9, "abcde")
	a4 := a.Delete(1, 2)
	fmt.Println("a", a)
	//fmt.Println(string(a.ToJSON()))
	fmt.Println("a1", a1)
	//fmt.Println(string(a1.ToJSON()))
	fmt.Println("a2", a2)
	//fmt.Println(string(a2.ToJSON()))
	fmt.Println("a3", a3)
	fmt.Println(string(a3.ToJSON()))
	fmt.Println("a4", a4.ToJSON())
	fmt.Println(a4.Report(1, 20))
}
