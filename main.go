package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/everettraven/wranglr/pkg/cmd"
)

func main() {
	if err := fang.Execute(context.TODO(), cmd.NewRootCommand()); err != nil {
		os.Exit(1)
	}
}
