//Package rope implements a persistent rope-like data structure.
//Persistent means that every operation does not modify the original
//objects.
//Refer to README.md for further information

package rope

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
)

//Rope represents a persistent rope data structure
type Rope struct {
	value  []rune
	weight int
	left   *Rope
	right  *Rope
}

func (rope *Rope) isLeaf() bool {
	return rope.left == nil
}

//New returns a new rope initialized with given string
func New(bootstrap string) *Rope {
	return &Rope{value: []rune(bootstrap), weight: utf8.RuneCountInString(bootstrap)}
}

//Len returns the length of the rope underlying string
func (rope *Rope) Len() int {
	if rope == nil {
		return 0
	}
	return rope.weight + rope.right.Len()
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
	Value  []rune
	Weight int
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
		Value:  rope.value,
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
	//Special case: if the other rope is nil, just return the first rope
	if other == nil {
		return rope
	}
	//Return a new rope with left and right subropes assigned respectively
	//to 'rope' and 'other'. Weight is the len of left rope.
	return &Rope{
		weight: rope.Len(),
		left:   rope,
		right:  other,
	}
}

//Internal function used by Split function.
//It accepts idx to split (1-based), a slice for the rope parts
//to be used for the second rope, a slice for the rope whose weight
//must be updated, and a slice to record weights to remove
func (rope *Rope) internalSplit(idx int,
	snd []*Rope,
	upd []*Rope,
	updCount []int) (*Rope, []*Rope, []*Rope, []int) {
	//If idx is equal to rope weight, we're arrived:
	//- Update upd slice with the first weight to remove from parents
	//- Add right rope as starter for the second rope slice;
	//- Create a rope equal to the original but without right rope
	if idx == rope.weight {
		updCount = append(updCount, rope.right.Len())
		return &Rope{
			weight: rope.weight,
			value:  rope.value,
			left:   rope.left,
		}, append(snd, rope.right), upd, updCount
	} else if idx > rope.weight {
		//We have to go right, call the function on right side with appropriate index.
		//Builds the rope and pass it up the stack with the other parameters.
		newRight, snd, upd, updCount := rope.right.internalSplit(idx-rope.weight, snd, upd, updCount)
		return &Rope{
			weight: rope.weight,
			value:  rope.value,
			left:   rope.left,
			right:  newRight,
		}, snd, upd, updCount
	} else if idx < rope.weight && !rope.isLeaf() {
		newLeft, snd, upd, updCount := rope.left.internalSplit(idx, snd, upd, updCount)
		snd = append(snd, rope.right)
		fst := &Rope{
			weight: rope.weight,
			left:   newLeft,
		}
		upd = append(upd, fst)
		updCount = append(updCount, rope.right.Len())
		return fst, snd, upd, updCount
	} else {
		sndW := len(rope.value) - idx
		updCount := append(updCount, sndW)
		snd = append(snd, &Rope{weight: sndW, value: rope.value[idx:len(rope.value)]})
		return &Rope{
			weight: idx,
			value:  rope.value[0:idx],
		}, snd, upd, updCount
	}
}

//Split generates two strings starting from one,
//splitting it at input index (1-based)
func (rope *Rope) Split(idx int) (firstRope *Rope, secondRope *Rope) {
	snd := make([]*Rope, 0)
	upd := make([]*Rope, 0)
	updCount := make([]int, 0)
	firstRope, snd, upd, updCount = rope.internalSplit(idx, snd, upd, updCount)
	for i, r := range snd {
		if r != nil {
			if i == 0 {
				secondRope = r
			} else {
				secondRope = secondRope.Concat(r)
			}
		}
	}
	var w int
	for i, u := range upd {
		w += updCount[i]
		u.weight = u.weight - w
	}
	return firstRope, secondRope
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
//for len bytes
func (rope *Rope) Report(idx int, lenght int) []rune {
	//if idx > rope.weight we go right with modified idx
	if idx > rope.weight {
		return rope.right.Report(idx-rope.weight, lenght)
	} else
	//if idx <= rope.weight we check if the left branch
	//has enough values to fetch report, else we split
	if rope.weight >= idx+lenght-1 {
		//we have enough space, just go left (if there is a left!)
		if !rope.isLeaf() {
			return rope.left.Report(idx, lenght)
		} else {
			//we're in a leaf, fetch from here
			return rope.value[idx-1 : idx+lenght-1]
		}
	} else {
		//Split the work and then merge both parts
		l := rope.left.Report(idx, rope.weight-idx+1)
		r := rope.right.Report(1, lenght-rope.weight+idx-1)
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
