//Package rope implements a persistent rope-like data structure.
//Persistent means that every operation does not modify the original
//objects.
//Refer to README.md for further information

package rope

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"unicode/utf8"
)

//Rope represents a persistent rope data structure
type Rope struct {
	value  []rune
	weight int
	length int
	left   *Rope
	right  *Rope
}

func (rope *Rope) isLeaf() bool {
	return rope.left == nil
}

//New returns a new rope initialized with given string
func New(bootstrap string) *Rope {
	return &Rope{
		value:  []rune(bootstrap),
		weight: utf8.RuneCountInString(bootstrap),
		length: len([]rune(bootstrap))}
}

//Len returns the length of the rope underlying string
func (rope *Rope) Len() int {
	if rope == nil {
		return 0
	}
	return rope.length
}

//String returns the complete string stored in the rope
func (rope *Rope) String() string {
	if rope.isLeaf() {
		return string(rope.value)
	}
	s1 := rope.left.String()
	if rope.right != nil {
		return s1 + rope.right.String()
	} else {
		return s1
	}
}

//Internal struct for generating JSON
type ropeForJSON struct {
	Value  string
	Weight int
	Length int
	Left   *ropeForJSON
	Right  *ropeForJSON
}

//Utility function that transforms a *Rope in a *ropeForJSON
func (rope *Rope) toRopeForJSON() *ropeForJSON {
	if rope == nil {
		return nil
	}
	return &ropeForJSON{
		Weight: rope.weight,
		Value:  string(rope.value),
		Length: rope.length,
		Left:   rope.left.toRopeForJSON(),
		Right:  rope.right.toRopeForJSON(),
	}
}

//ToJSON generates a indented JSON rope conversion
func (rope *Rope) ToJSON() string {
	rope2 := rope.toRopeForJSON()
	var out bytes.Buffer
	b, _ := json.Marshal(rope2)
	json.Indent(&out, b, "", "  ")
	return string(out.Bytes())
}

//Index retrieves the byte at rope position idx (1-based)
func (rope *Rope) Index(idx int) rune {
	if idx > rope.weight {
		return rope.right.Index(idx - rope.weight)
	} else if rope.isLeaf() {
		return rope.value[idx-1]
	} else {
		return rope.left.Index(idx)
	}
}

//Concat merges two ropes and generates a brand new one
func (rope *Rope) Concat(other *Rope) *Rope {
	//Special case: if the first rope is nil, just return the second rope
	if rope == nil {
		return other
	}
	//Special case: if the other rope is nil, just return the first rope
	if other == nil {
		return rope
	}
	//Return a new rope with 'rope' and 'other' assigned respectively
	//to left and right subropes. Weight is the len of left rope.
	return &Rope{
		weight: rope.Len(),
		length: rope.length + other.length,
		left:   rope,
		right:  other,
	}
}

//Internal function used by Split function.
//It accepts idx to split (1-based), a slice for the rope parts
//to be used for the second rope, a slice for the rope whose weight
//must be updated, and a slice to record weights to remove
func (rope *Rope) split(idx int,
	secondRope *Rope) (*Rope, *Rope) {
	//If idx is equal to rope weight, we're arrived:
	//- If is leaf, return it;
	//- Otherwise, return its left rope.
	//Right rope initialises secondRope.
	if idx == rope.weight {
		var r *Rope
		if rope.isLeaf() {
			r = rope
		} else {
			r = rope.left
		}
		return r, rope.right
	} else if idx > rope.weight {
		//We have to go right, call the function on right side with appropriate index.
		//Builds the rope and pass it up the stack with the other parameters.
		newRight, secondRope := rope.right.split(idx-rope.weight, secondRope)
		return &Rope{
			weight: rope.weight,
			left:   rope.left,
			right:  newRight,
		}, secondRope
	} else if rope.isLeaf() { //idx < rope.weight, go left!
		//It's a leaf: we have to create a new rope by splitting leaf at index
		return &Rope{
				weight: idx,
				value:  rope.value[0:idx],
				length: idx,
			}, secondRope.Concat(&Rope{
				weight: len(rope.value) - idx,
				value:  rope.value[idx:len(rope.value)],
				length: len(rope.value) - idx,
			})
	} else {
		newLeft, secondRope := rope.left.split(idx, secondRope)
		return newLeft, secondRope.Concat(rope.right)
	}
}

//Split generates two strings starting from one,
//splitting it at input index (1-based)
func (rope *Rope) Split(idx int) (firstRope *Rope, secondRope *Rope) {
	//Create the slices for split
	return rope.split(idx, secondRope)
}

//Insert generates a new rope inserting a string into the
//original rope
func (rope *Rope) Insert(idx int, s string) *Rope {
	//Split rope at insert point
	r1, r2 := rope.Split(idx)
	//Rejoin the two split parts with string to insert as middle rope
	return r1.Concat(New(s)).Concat(r2)
}

//Delete generates a new rope by deleting len characters
//from the original one starting at character idx
func (rope *Rope) Delete(idx int, len int) *Rope {
	r1, r2 := rope.Split(idx - 1)
	_, r4 := r2.Split(len)
	return r1.Concat(r4)
}

//Report return a substring of the rope starting from idx included (1-based)
//for len runes
func (rope *Rope) Report(idx int, length int) string {
	return string(rope.internalReport(idx, length))
}

func (rope *Rope) internalReport(idx int, length int) []rune {
	//if idx > rope.weight we go right with modified idx
	if idx > rope.weight {
		return rope.right.internalReport(idx-rope.weight, length)
	} else
	//if idx <= rope.weight we check if the left branch
	//has enough values to fetch report, else we split
	if rope.weight >= idx+length-1 {
		//we have enough space, just go left (if there is a left!)
		if !rope.isLeaf() {
			return rope.left.internalReport(idx, length)
		} else {
			//we're in a leaf, fetch from here
			return rope.value[idx-1 : idx+length-1]
		}
	} else {
		//Split the work and then merge both parts
		l := rope.left.internalReport(idx, rope.weight-idx+1)
		r := rope.right.internalReport(1, length-rope.weight+idx-1)
		s := make([]rune, len(l)+len(r))
		for i := 0; i < len(l); i++ {
			s[i] = l[i]
		}

		for i := 0; i < len(r); i++ {
			s[i+len(l)] = r[i]
		}

		return s

	}
}
