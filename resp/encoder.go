package resp

import (
	"fmt"
)

func Encode(data any) []byte {

	switch data := data.(type) {
	case string:
		return encodeString(data)
	case error:
		return encodeError(data)
	case int:
		return encodeInt(data)
	case nil:
		return encodeBulkString(nil)
	default:
		return encodeError(fmt.Errorf("unknown response type %T", data))
	}
}

func EncodeArray(data ...any) []byte {
	res := make([]byte, 0)
	header := fmt.Sprintf("%c%d\r\n", SymbolArray, len(data))
	res = append(res, []byte(header)...)

	for _, entry := range data {
		val := Encode(entry)
		res = append(res, val...)
	}

	return res
}

func EncodeBulkString(data string) []byte {
	return encodeBulkString(&data)
}

func encodeString(data string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", data))
}

func encodeBulkString(data *string) []byte {
	if data == nil {
		return []byte("$-1\r\n")
	}

	return []byte(fmt.Sprintf("%c%d\r\n%s\r\n", SymbolBulkString, len(*data), *data))

}

func encodeError(err error) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", err.Error()))
}

func encodeInt(data int) []byte {

	return []byte(fmt.Sprintf("%c%d\r\n", SymbolInt, data))
}
