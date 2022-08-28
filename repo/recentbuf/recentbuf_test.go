package recentbuf

import (
	"testing"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"reflect"
	"encoding/json"
	"strings"
	"fmt"
)

func TestAddInitial(t *testing.T) {

	buf := New()

	for i := 0; i < MaxSize; i++ {
		ID := id.New()
		buf.Add(ID)

		expected := i + 1
		if len(buf.elements) != expected {
			t.Errorf("Buffer has %d elements, want %d", len(buf.elements), expected)
		}

	}

}

func TestFull(t *testing.T) {

	buf := New()

	// fill the buffer
	for i := 0; i < MaxSize; i++ {
		ID := id.New()
		buf.Add(ID)

	}

	testCount := 3
	expectHead := make([]id.ID, 0, testCount)
	for i := 1; i <= testCount; i++ {

		// add one more ID
		ID := id.New()
		buf.Add(ID)

		expectHead = append(expectHead, id.Zero)
		copy(expectHead[1:], expectHead[:i])
		expectHead[0]=ID


		//check
		if len(buf.elements) != MaxSize {
			t.Errorf("buffer size is %d, want %d", len(buf.elements), MaxSize)
		}

		head := buf.elements[:i]

		if !reflect.DeepEqual(head, expectHead) {
			t.Errorf("After pushing %d element to the full buffer, the head of the buffer is %s, want %s", i, head, expectHead)
		}

	}

}

func TestBuffer_Add_Duplicate(t *testing.T) {

	id1 := id.New()
	id2 := id.New()
	id3 := id.New()

	idNames:=map[id.ID]string{id1:"1",id2:"2",id3:"3"}

	IDsToNames:=func(IDs []id.ID)[]string{ // return slice of short names for the ids
		result:=make([]string, len(IDs))

		for i:=range IDs{
			result[i]=fmt.Sprintf("( id: %s)",idNames[IDs[i]])
		}
		return result
	}

	tests := []struct {
		Case string  // test case name
		Seq  []id.ID // sequence of ids to add
		Exp  []id.ID // expected elements
	}{
		{Case:"two differrent, newer must have lesser index", Seq:[]id.ID{id1, id2}, Exp:[]id.ID{id2, id1}},
		{Case:"three differrent, newer must have lesser index", Seq:[]id.ID{id1, id2, id3}, Exp:[]id.ID{id3, id2, id1}},
		{Case:"two identical", Seq:[]id.ID{id1, id1}, Exp:[]id.ID{id1}},
		{Case:"three identical", Seq:[]id.ID{id1, id1, id1}, Exp:[]id.ID{id1}},
		{Case:"one different, two identical, newer must have lesser index", Seq:[]id.ID{id2, id1, id1}, Exp:[]id.ID{id1, id2}},
		{Case:"two identical, one different, newer must have lesser index", Seq:[]id.ID{id1, id1, id2}, Exp:[]id.ID{id2, id1}},
		{Case:"duplicate must go to the head", Seq:[]id.ID{id1, id3, id2, id3}, Exp:[]id.ID{id3, id2, id1}},
		{Case:"duplicate must go to the head", Seq:[]id.ID{id1, id3, id2, id3, id3}, Exp:[]id.ID{id3, id2, id1}},
	}

	for _, test := range tests {
		buf := New()
		for _, ID := range test.Seq {
			buf.Add(ID)
		}
		result := buf.Elements()
		expected := test.Exp
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Test case %s:\nafter adding sequence \t%s,\nElements() are\t\t\t%s,\nwant\t\t\t\t\t%s",
				strings.ToUpper(test.Case),IDsToNames(test.Seq), IDsToNames(result), IDsToNames(expected))
		}
	}
}

func TestJSON(t *testing.T) {

	buf := New()

	bufLen := 3
	for i := 0; i < bufLen; i++ {

		buf.Add(id.New())
	}

	data, err := json.Marshal(buf)

	if err != nil {
		t.Error(err)
	}

	buf2 := &Buffer{}
	if err := json.Unmarshal(data, buf2); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(buf.elements, buf2.elements) {
		t.Errorf("after json marshal-umarshal reoundtrip returned buffer is %s, want %s", buf2, buf)
	}

}
