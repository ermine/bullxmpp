package jid
import "testing"
import "fmt"
import "os"

func TestIDN(t *testing.T) {
	// var str = "<html><body><font size=\"48\">i&apos;lold%20you!</font><body><html>@xn----7sbcgldb2a9a7ach1k.xn--p1ai"
	var str = "xn----7sbcgldb2a9a7ach1k.xn--p1ai"
	var jid, err = New(str)
	if err != nil {
		t.Fatal(err)
	}
	var ret, err1 = jid.GetIDN()
	if err1 != nil {
		t.Fatal(err1)
	}
	fmt.Println(ret)
}
