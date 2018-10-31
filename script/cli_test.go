package script

import "testing"

func TestToSliceStringMap(t *testing.T) {
	slice := make([]interface{}, 2)
	slice[0] = make(map[string]interface{}, 0)
	slice[1] = "foo"

	sliceMap := ToSliceStringMap(slice)
	if len(sliceMap) == 2 {
		t.Fatal("got ok for slice with string")
	}

	slice[1] = make(map[string]interface{}, 0)
	sliceMap = ToSliceStringMap(slice)
	if len(sliceMap) != 2 {
		t.Fatal("did not get ok for slice of map[string]interface{}")
	}
}
