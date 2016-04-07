package ifitlib

import "testing"

/*

TODO:

*/

func Test_Exists(t *testing.T) {
	b := Exists("lib.go")
	if !b {
		t.Errorf("Exists: failed to find file\n")
	}
	b = Exists("noexistent.go")
	if b {
		t.Errorf("Exists: found file that is nonexistent\n")
	}
}

func Test_InArray(t *testing.T) {
	b := InArray("aa", []string{"bb", "aa", "cc"})
	if !b {
		t.Errorf("InArray: failed to find\n")
	}
	b = InArray("aa", []string{"bb", "AA", "cc"})
	if b {
		t.Errorf("InArray: found someting that was not there\n")
	}
}

func Test_Json1(t *testing.T) {

	mm, err := JsonStringToString(`{"aa":"bb"}`)
	if err != nil {
		t.Errorf("JsonStringToString: got parse error %s\n", err)
	}
	bv, b := mm["aa"]
	if !b {
		t.Errorf("JsonStringToString: failed to find key aa\n")
	}
	if bv != "bb" {
		t.Errorf("JsonStringToString: incorrect value\n")
	}
}

func Test_Json2(t *testing.T) {
	mm, err := JsonStringToStringString(`{"aa":{"cc":"bb"}}`)
	if err != nil {
		t.Errorf("JsonStringToStringString: got parse error %s\n", err)
	}
	t1, b := mm["aa"]
	if !b {
		t.Errorf("JsonStringToStringString: failed to find key aa\n")
	}
	bv, b := t1["cc"]
	if !b {
		t.Errorf("JsonStringToStringString: failed to find key cc\n")
	}
	if bv != "bb" {
		t.Errorf("JsonStringToStringString: incorrect value\n")
	}

}

func Test_ParseLineIntoWords(t *testing.T) {
	// func ParseLineIntoWords(line string) []string {
	ss := ParseLineIntoWords("a b c")
	if len(ss) != 3 {
		t.Errorf("ParseLineIntoWords: expected 3 resutls got %d\n", len(ss))
	}
	if ss[0] != "a" {
		t.Errorf("ParseLineIntoWords: expected 'a' resutls got %s\n", ss[0])
	}
	if ss[1] != "b" {
		t.Errorf("ParseLineIntoWords: expected 'b' resutls got %s\n", ss[1])
	}
	if ss[2] != "c" {
		t.Errorf("ParseLineIntoWords: expected 'c' resutls got %s\n", ss[2])
	}

	vv := ParseLineIntoWords(`a "b bb bbb" c`)
	if len(vv) != 3 {
		t.Errorf("ParseLineIntoWords: expected 3 resutls got %d\n", len(vv))
	}
	if vv[0] != "a" {
		t.Errorf("ParseLineIntoWords: expected 'a' resutls got %s\n", vv[0])
	}
	if vv[1] != "\"b bb bbb\"" {
		t.Errorf("ParseLineIntoWords: expected \"b bb bbb\" resutls got %s\n", vv[1])
	}
	if vv[2] != "c" {
		t.Errorf("ParseLineIntoWords: expected 'c' resutls got %s\n", vv[2])
	}
}

func Test_KeysSorted(t *testing.T) {
	//	func KeysSorted(sub map[string]string) (strs []string) {
	aMap := make(map[string]string)
	aMap["abc"] = "abc"
	aMap["def"] = "123"
	aMap["111"] = "yep"
	aMap["222"] = "nope"

	mm := KeysSorted(aMap)
	if len(mm) != 4 {
		t.Errorf("KeysSorted: expected 4 resutls got %d\n", len(mm))
	}
	ref := []string{"111", "222", "abc", "def"}
	for i := 0; i < 4; i++ {
		if len(mm) != 4 {
			t.Errorf("KeysSorted: expected %s at %d got %s\n", ref[i], i, mm[i])
		}
	}
}

func Test_CommaList(t *testing.T) {
	//	func CommaList(strs []string) (s string) {
	ss := CommaList([]string{"aa", "bb", "cc"})
	if ss != "aa, bb, cc" {
		t.Errorf("CommaList: expected 'aa, bb, cc' resutls got %s\n", ss)
	}
	ss = CommaList([]string{})
	if ss != "" {
		t.Errorf("CommaList: expected '' resutls got %s\n", ss)
	}
	ss = CommaList([]string{"aa"})
	if ss != "aa" {
		t.Errorf("CommaList: expected 'aa' resutls got %s\n", ss)
	}
}

func Test_ParseNameValueOpt(t *testing.T) {
	//	func ParseNameValueOpt(s string) (name, value string, err error) {

	name, value, err := ParseNameValueOpt("abc=def")
	if err != nil {
		t.Errorf("ParseNameValueOpt: got error when not expected\n")
	}
	if name != "abc" {
		t.Errorf("ParseNameValueOpt: expected 'abc' resutls got %s\n", name)
	}
	if value != "def" {
		t.Errorf("ParseNameValueOpt: expected 'def' resutls got %s\n", value)
	}

	name, value, err = ParseNameValueOpt("abc")
	if err != nil {
		t.Errorf("ParseNameValueOpt: got error when not expected\n")
	}
	if name != "abc" {
		t.Errorf("ParseNameValueOpt: expected 'abc' resutls got %s\n", name)
	}
	if value != "on" {
		t.Errorf("ParseNameValueOpt: expected 'on' resutls got %s\n", value)
	}

}

func Test_FindFile(t *testing.T) {
	//	func FindFile(PathOfInput, fn string, sp []string) (rv string) {
	// TODO - should use array of tests and a more comprehensive test set for this
	rv := FindFile("./", "t.1", []string{"t1", "./t2", "../ifitlib/t3"})
	//fmt.Printf("rv=[%s]\n", rv)
	if rv != "../ifitlib/t3/t.1" {
		t.Errorf("FindFile: expected '../ifitlib/t3/t.1' resutls got %s\n", rv)
	}
	rv = FindFile("./", "t.1", []string{"t1", "./t2", "./t3"})
	//fmt.Printf("rv=[%s]\n", rv)
	if rv != "t3/t.1" {
		t.Errorf("FindFile: expected 't3/t.1' resutls got %s\n", rv)
	}
	rv = FindFile("", "t.1", []string{"./t1", "./t2", "./t3"})
	//fmt.Printf("rv=[%s]\n", rv)
	if rv != "t3/t.1" {
		t.Errorf("FindFile: expected 't3/t.1' resutls got %s\n", rv)
	}
}

//	func GetItemSet(line string, nthItem int) (set []string) {
func Test_GetItemN_and_GetItemSet(t *testing.T) {
	//	func GetItemN(line string, nthItem int) (name string) {
	name := GetItemN("a b c e f", 3)
	if name != "c" {
		t.Errorf("GetItemN: expected 'c' resutls got %s\n", name)
	}
	name = GetItemN("a b", 3)
	if name != "" {
		t.Errorf("GetItemN: expected '' resutls got %s\n", name)
	}
	name = GetItemN("a b", -3)
	if name != "" {
		t.Errorf("GetItemN: expected '' resutls got %s\n", name)
	}
	ss := GetItemSet("a b c e !!", 3)
	if len(ss) != 2 {
		t.Errorf("GetItemSet: expected 2 resutls got %d\n", len(ss))
	}
	ss = GetItemSet("a b c e f", 3)
	if len(ss) != 3 {
		t.Errorf("GetItemSet: expected 3 resutls got %d\n", len(ss))
	}
	ss = GetItemSet("a b", 3)
	if len(ss) != 0 {
		t.Errorf("GetItemSet: expected 0 resutls got %d\n", len(ss))
	}
	ss = GetItemSet("a b c e f", -3)
	if len(ss) != 3 {
		t.Errorf("GetItemSet: expected 3 resutls got %d\n", len(ss))
	}
}
