package printers

import (
	"fmt"
	"strings"

	"github.com/everettraven/synkr/pkg/plugins"
	"github.com/everettraven/synkr/pkg/plugins/github"
)

type Markdown struct{}

func (md *Markdown) Print(results ...plugins.SourceResult) error {
	for _, result := range results {
		out, err := md.PrintResult(result)
		if err != nil {
			return err
		}
		fmt.Println(out)
	}

	return nil
}

func (md *Markdown) PrintResult(result plugins.SourceResult) (string, error) {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("# %s - %s\n", result.Source, result.Project))

	for _, item := range result.Items {
		switch item := item.(type) {
		case github.RepoItem:
			out.WriteString(md.PrintGitHubRepoItem(item))
		default:
			out.WriteString(fmt.Sprintf("%v", item))
		}

		out.WriteString("\n")
	}

	return out.String(), nil
}

func (md *Markdown) PrintGitHubRepoItem(item github.RepoItem) string {
	var out strings.Builder

	out.WriteString(fmt.Sprintf("## [%s][%s]: %s \n", item.Type, item.State, item.Title))
	out.WriteString(fmt.Sprintf("**URL**: %s\n", item.URL))
	out.WriteString(fmt.Sprintf("**Author**: *%s*\n", item.Author))
	out.WriteString(fmt.Sprintf("**Assignees**: %s\n", strings.Join(item.Assignees, ",")))

	out.WriteString("\n")
	labelStrs := []string{}
	for _, label := range item.Labels {
		labelStrs = append(labelStrs, fmt.Sprintf("`%s`", label))
	}
	out.WriteString(fmt.Sprintf("%s\n\n", strings.Join(labelStrs, " ")))
	out.WriteString(fmt.Sprintf("%s\n", item.Body))

	return out.String()
}
