package printers

import (
	"encoding/json"
	"fmt"

	"github.com/everettraven/synkr/pkg/engine"
)

type JSON struct{}

func (j *JSON) Print(results ...engine.SourceResult) error {
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
