package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'
)

type Value struct {
	Typ   byte
	Str   string
	Bulk  string
	Array []Value
}

type Reader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}

func (r *Reader) ReadLine() (line []byte, err error){
	
	line, err = r.reader.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	return line[:len(line) - 2], nil
}

func (r *Reader) ReadInteger()(int, error){

	line, err := r.ReadLine()

	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (r *Reader) ReadValue() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return r.parseArray()
	case BULK_STRING:
		return r.parseBulkString()
	case SIMPLE_STRING:
		return r.parseSimpleString()
	case ERROR:
		return r.parseError()
	default:
		return Value{}, fmt.Errorf("unsupported RESP type: %c", typ)
	}
}

func (r *Reader) parseArray() (Value, error) {

	len, err := r.ReadInteger()	

	if err != nil {
		return Value{}, err
	}

	array := make([]Value,  len)

	for i := 0; i < len; i++{
		val, err := r.ReadValue()
		if err != nil {
			return Value{}, err
		}

		array[i] =  val
	}

	return Value{Typ: ARRAY, Array: array}, nil

}


	func (r *Reader) parseBulkString() (Value, error) {

	len, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	buf := make([]byte, len+2) 

	if _, err := io.ReadFull(r.reader, buf); err != nil {
		return Value{}, err
	}

	return Value{Typ: BULK_STRING, Bulk: string(buf[:len])}, nil
}


func (r *Reader) parseSimpleString()(Value, error){

	line, err := r.ReadLine()

	if err != nil{
		return Value{}, err
	}

	return Value{Typ: SIMPLE_STRING, Str: string(line)}, nil

}

func (r *Reader) parseError() (Value, error) {
	line, err := r.ReadLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: ERROR, Str: string(line)}, nil
}