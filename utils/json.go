package utils

import (
	"os"
	"github.com/bytedance/sonic"
)

// JsonBytes marshals v to JSON and returns its bytes represenation.
func JsonBytes[V any](v *V) ([]byte, error) {
	return sonic.Marshal(v)
}

// JsonBytes marshals v to JSON and returns its string represenation.
func JsonStr[V any](v *V) (string, error) {
	return sonic.MarshalString(v)
}

func JsonToFile[V any, D []byte | string](path string, perms os.FileMode, v *V, marshal func(v *V) (D, error), write func(path string, perms os.FileMode, data D) error) error {
	d, err := marshal(v)

	if err != nil {
		return err
	}

	if err := write(path, perms, d); err != nil {
		return err
	}

	return nil
}

// JsonBytesToFile marshals v to JSON and saves resulting bytes represenation into the file under path specified. Non-existent parent directories are recursively created.
func JsonBytesToFile[V any](path string, perms os.FileMode, v *V) error {
	return JsonToFile(path, perms, v, JsonBytes, WriteFileBytes)
}

// JsonStrToFile marshals v to JSON and saves resulting string represenation into the file under path specified. Non-existent parent directories are recursively created.
func JsonStrToFile[V any](path string, perms os.FileMode, v *V) error {
	return JsonToFile(path, perms, v, JsonStr, WriteFileString)
}

// JsonParseBytes parses expected JSON bytes into the target T data structure under pointer.
func JsonParseBytes[T any](source []byte, target *T) error {
	return sonic.Unmarshal(source, target)
}

// JsonParseStr parses expected JSON string into the target T data structure under pointer.
func JsonParseStr[T any](source string, target *T) error {
	return sonic.UnmarshalString(source, target)
}

// JsonParseBytesFromFile loads and parses expected JSON bytes into the target T data structure under pointer.
func JsonParseBytesFromFile[T any](path string, target *T) error {
	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	return JsonParseBytes(data, target)
}
