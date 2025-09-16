package printers

import (
	"encoding/json"
	"fmt"
	"go.starlark.net/starlark"
)

type JSON struct{}

func (j *JSON) Print(results ...starlark.Value) error {
	outBytes := []byte{}
	for _, result := range results {
		out, err := json.Marshal(result)
		if err != nil {
			return err
		}
		outBytes = append(outBytes, out...)
	}

	fmt.Println(string(outBytes))
	return nil
}
