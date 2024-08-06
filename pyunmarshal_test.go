package pyunmarshal

import "testing"

var testDataEmpty = []byte{}

var testData = []byte{}

func TestBufferEmpty(t *testing.T) {
	pyMarshal := NewMarshal(testDataEmpty)
	_, eof, _ := pyMarshal.Read()
	if !eof {
		t.Error("eof expected")
	}
}

func TestBufferSimple(t *testing.T) {
	pyMarshal := NewMarshal(testData)
	_, eof, _ := pyMarshal.Read()
	if !eof {
		t.Error("eof expected")
	}
}
