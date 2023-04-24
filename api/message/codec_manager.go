package message

import "fmt"

var registry = map[string]Codec{}

func RegisterCodec(name string, codec Codec) {
	registry[name] = codec
}

func GetCodec(name string) (Codec, error) {
	var err error
	codec, ok := registry[name]
	if !ok {
		err = fmt.Errorf("no codec registered with name '%s'", name)
	}
	return codec, err
}
