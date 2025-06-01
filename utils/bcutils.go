package utils

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/mr-tron/base58"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)

	if len(encoded) == 0 {
		log.Panic("Failed to encode input to Base58")
	}

	return []byte(encoded)
}

func Base58Decode(input []byte) []byte {
	decoded, err := base58.Decode(string(input[:]))

	if err != nil {
		log.Panicf("Failed to decode Base58 input: %v", err)
	}

	if len(decoded) == 0 {
		log.Panic("Decoded Base58 input is empty")
	}
	return decoded
}

func Serialize[T any](data T) []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(data)
	HandleError(err)
	return encoded.Bytes()
}

func Deserialize[T any](data []byte) T {
	var result T
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&result)
	HandleError(err)
	return result
}
