package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/everettraven/synkr/pkg/cmd"
)

func main() {
	if err := fang.Execute(context.TODO(), cmd.NewSynkrCommand()); err != nil {
		os.Exit(1)
	}
}
