package printers

import (
	"encoding/json"

	"github.com/everettraven/synkr/pkg/engine"
)

type JSON struct{}

func (j *JSON) Print(result engine.SourceResult) (string, error) {
	outBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(outBytes), nil
}
