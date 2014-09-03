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
	length int
	left   *Rope
	right  *Rope
}

//isLeaf returns true if the rope is a leaf
func (rope *Rope) isLeaf() bool {
	return rope.left == nil
}

//New returns a new rope initialized with given string
func New(bootstrap string) *Rope {
	len := utf8.RuneCountInString(bootstrap)
	return &Rope{
		value:  []rune(bootstrap),
		weight: len,
		length: len}
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
	return rope.Report(1, rope.length)
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

//Index retrieves the rune at rope position idx (1-based)
func (rope *Rope) Index(idx int) rune {
	if rope.isLeaf() {
		return rope.value[idx-1]
	} else if idx > rope.weight {
		return rope.right.Index(idx - rope.weight)
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
	//to left and right subropes.
	return &Rope{
		weight: rope.Len(),
		length: rope.Len() + other.Len(),
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
		//We have to recurse on right side.
		newRight, secondRope := rope.right.split(idx-rope.weight, secondRope)
		return &Rope{
			weight: rope.weight,
			left:   rope.left,
			right:  newRight,
		}, secondRope
	} else {
		//idx < rope.weight, we recurse on the left side
		if rope.isLeaf() {
			//It's a leaf: we have to create a new rope by splitting leaf at index
			var lr *Rope
			if idx > 0 {
				lr = &Rope{
					weight: idx,
					value:  rope.value[0:idx],
					length: idx,
				}
			}
			secondRope = &Rope{
				weight: len(rope.value) - idx,
				value:  rope.value[idx:len(rope.value)],
				length: len(rope.value) - idx,
			}
			return lr, secondRope
		} else {
			newLeft, secondRope := rope.left.split(idx, secondRope)
			return newLeft, secondRope.Concat(rope.right)
		}
	}
}

//Split generates two strings starting from one,
//splitting it at input index (1-based)
func (rope *Rope) Split(idx int) (firstRope *Rope, secondRope *Rope) {
	if rope == nil {
		return nil, nil
	}
	if idx <= 0 {
		return rope, nil
	}
	if idx >= rope.length {
		return nil, rope
	}
	//Create the slices for split
	return rope.split(idx, secondRope)
}

//Insert generates a new rope inserting a string into the
//original rope.
func (rope *Rope) Insert(idx int, str string) *Rope {
	if rope == nil {
		return New(str)
	}
	if idx < 0 {
		rope.Insert(0, str)
	}
	if idx > rope.length {
		rope.Insert(rope.length, str)
	}
	//Split rope at insert point
	r1, r2 := rope.Split(idx)
	//Rejoin the two split parts with string to insert as middle rope
	return r1.Concat(New(str)).Concat(r2)
}

//Delete generates a new rope by deleting len characters
//from the original one starting at character idx
func (rope *Rope) Delete(idx int, length int) *Rope {
	r1, r2 := rope.Split(idx - 1)
	_, r4 := r2.Split(length)
	return r1.Concat(r4)
}

//Report return a substring of the rope starting from idx included (1-based)
func (rope *Rope) Report(idx int, length int) string {
	if rope == nil || idx > rope.length || length < 1 {
		return ""
	}
	if idx < 1 {
		rope.Report(1, length)
	}
	if idx+length-1 > rope.length {
		rope.Report(idx, rope.length-idx+1)
	}
	res := make([]rune, length)
	rope.internalReport(idx, length, res)
	return string(res)
}

func (rope *Rope) internalReport(idx int, length int, res []rune) {
	//if idx > rope.weight we go right with modified idx
	if idx > rope.weight {
		rope.right.internalReport(idx-rope.weight, length, res)
	} else
	//if idx <= rope.weight we check if the left branch
	//has enough values to fetch report, else we split
	if rope.weight >= idx+length-1 {
		//we have enough space, just go left or take the string
		if rope.isLeaf() {
			//we're in a leaf, fetch from here
			copy(res, rope.value[idx-1:idx+length-1])
		} else {
			rope.left.internalReport(idx, length, res)
		}
	} else {
		//Split the work
		rope.left.internalReport(idx, rope.weight-idx+1, res[:rope.weight])
		rope.right.internalReport(1, length-rope.weight+idx-1, res[rope.weight:])
	}
}

func (rope *Rope) Substr(idx int, length int) *Rope {
	if rope == nil || idx > rope.length || length < 1 {
		return nil
	}
	if idx < 1 {
		rope.Report(1, length)
	}
	if idx+length-1 > rope.length {
		rope.Report(idx, rope.length-idx+1)
	}

	_, r1 := rope.Split(idx - 1)
	r2, _ := r1.Split(length)
	return r2
}
