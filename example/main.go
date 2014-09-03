//Example of usage
package main

import (
	"fmt"
	"github.com/vinzmay/go-rope"
)

func main3() {
	a := rope.New("Hello, 世界")
	fmt.Println("a:", a)
	fmt.Println("a len:", a.Len())
	a1, a2 := a.Split(7)
	fmt.Println("a1:", a1.ToJSON())
	fmt.Println("a2:", a2.ToJSON())
	fmt.Println("a_test:", a1.Concat(a2))
	fmt.Println("a[8]:", string(a.Index(8)))
	fmt.Println("a[8-9]:", string(a.Report(8, 2)))
	main2()
}

func main() {
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
	a1, a2 := a.Split(13)
	//fmt.Println("a", a.ToJSON())
	fmt.Println("a1", a1)
	fmt.Println("a1j", a1.ToJSON())
	fmt.Println("a2", a2)
	fmt.Println("a2j", a2.ToJSON())
	/*fmt.Println("a", a)
	fmt.Println(a.Report(2, 4))
	fmt.Println(a.Substr(2, 4))*/
}

func main2() {
	a := rope.New("abc")
	a1, a2 := a.Split(1)
	fmt.Println(a)
	fmt.Println(a1)
	fmt.Println(a2)
}
