package pyunmarshal

import (
	"encoding/binary"
	"fmt"
)

const (
	DataTypeNil    = '0'
	DataTypeString = 's'
	DataTypeUTF8   = 'u'
	DataTypeDict   = '{'
	DataTypeTrue   = 'T'
	DataTypeFalse  = 'F'
	DataTypeInt    = 'i'
)

type pyMarshal struct {
	buffer []byte
	cursor uint
}

func NewMarshal(buffer []byte) pyMarshal {
	marshal := pyMarshal{buffer, 0}
	return marshal
}

// Current position in the buffer
func (t pyMarshal) pos() uint {
	return t.cursor
}

// Length of buffer
func (t pyMarshal) len() uint {
	return uint(len(t.buffer))
}

// Returns true if there is more data left to read through in the buffer
func (t pyMarshal) moreData() bool {
	return (t.pos() < t.len())
}

// Return current buffer position and increments it by `byteCount`
func (t *pyMarshal) advance(byteCount uint) uint {
	start := t.cursor
	t.cursor += byteCount
	return start
}

func (t *pyMarshal) readType() rune {
	return rune(t.readByte())
}

func (t *pyMarshal) readByte() byte {
	return t.buffer[t.advance(1)]
}

func (t *pyMarshal) readInt32() int {
	return int(t.readUInt32())
}

func (t *pyMarshal) readUInt32() uint {
	start := t.advance(4)
	b := t.buffer[start:t.pos()]
	return uint(binary.LittleEndian.Uint32(b))
}

// Read the buffer at the current location for `len` bytes and return the bytes read
func (t *pyMarshal) readBuffer(len uint) []byte {
	start := t.advance(len)
	return t.buffer[start:t.pos()]
}

// Reads the buffer at the current position and returns a deserialized object.
// If the buffer is not fully Read, then extra calls can be made to extract more data.
func (t *pyMarshal) Read() (object any, eof bool, err error) {
	// No more data left in the tank
	if !t.moreData() {
		return nil, true, nil
	}

	// Read the type first so we know what we are dealing with
	anchorIdx := t.pos()
	pyType := t.readType()
	switch pyType {
	case DataTypeNil:
		return nil, false, nil

	case DataTypeTrue:
		return true, false, nil

	case DataTypeFalse:
		return false, false, nil

	case DataTypeUTF8:
		fallthrough
	case DataTypeString:
		strLen := t.readUInt32()
		strBytes := t.readBuffer(strLen)
		return string(strBytes), false, nil

	case DataTypeInt:
		return t.readInt32(), false, nil

	case DataTypeDict:
		d := map[string]any{}
		for {
			key, end, err := t.Read()
			if err != nil {
				return nil, end, err
			}

			// Break out of loop if we didn't find a string for a key
			dictDone := false
			switch key.(type) {
			case string:
				break
			default:
				dictDone = true
			}

			if dictDone {
				break
			}

			// Get value
			value, end, err := t.Read()
			if err != nil {
				return nil, end, err
			}

			d[key.(string)] = value
		}
		return d, false, nil

	default:
		return nil, false, fmt.Errorf("unsupported type 0x%X at byte 0x%X", pyType, anchorIdx)
	}
}
