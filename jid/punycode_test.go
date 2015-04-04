package jid
import "testing"
import "fmt"

func TestEncode(t *testing.T) {
	var str = "--7sbcgldb2a9a7ach1k"
	var ret, err = punycode_decode([]rune(str))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("decoded " + string(ret))

	ret, err = punycode_encode(ret)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("encoded [" + string(ret) + "]")

	if str != string(ret) {
		t.Fatal("bad" + string(ret))
	}
}
