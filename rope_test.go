package rope

import "testing"

//Test rope creation
func TestRopeCreation(t *testing.T) {
	r := New("test")
	if r.String() != "test" {
		t.Error("Error creating rope - equality fail: ", r, " != test")
	}
	if r.Len() != 4 {
		t.Error("Error creating rope - length fail: ", r.Len(), "!= 4")
	}
}

//Test rope concatenation
func TestRopeConcat(t *testing.T) {
	r := New("abcdef")
	r2 := New("ghilmno")
	r3 := r.Concat(r2)
	if r.String() != "abcdef" || r.Len() != 6 {
		t.Error("Error concatenating ropes, r modified:", r)
	}
	if r2.String() != "ghilmno" || r2.Len() != 7 {
		t.Error("Error concatenating ropes, r2 modified:", r2)
	}
	if r3.String() != "abcdefghilmno" || r3.Len() != 13 {
		t.Error("Error concatenating ropes, r3 not correct:", r3, "!= abcdefghilmno")
	}
}

//Test rope split
func TestRopeSplit(t *testing.T) {
	r := New("abcdef")
	r1, r2 := r.Split(4)
	if r.String() != "abcdef" || r1.String() != "abcd" || r2.String() != "ef" {
		t.Error("Error splitting string: abcd/ef => ", r1, r2)
	}
}
