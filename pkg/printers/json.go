package printers

import (
	"encoding/json"
	"fmt"

	"github.com/everettraven/synkr/pkg/plugins"
)

type JSON struct{}

func (j *JSON) Print(results ...plugins.SourceEntry) error {
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
