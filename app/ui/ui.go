//go:generate mockgen -package=ui -source=ui.go -destination=ui_mock.go

package ui

import (
	"io"

	"fmt"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
)

// UI is an interface for intaracting with the shell
type UI interface {
	Ask(query string) (string, error)
	Error(msg string)
	Status(prefix, message string, colorAttrs ColorAttrs)
}

type ui struct {
	cli.Ui
}

// NewUI creates a new UI instance from streams
func NewUI(inStream io.Reader, outStream, errStream io.Writer) UI {
	return &ui{
		Ui: &cli.PrefixedUi{
			ErrorPrefix: uiPrefix("Error", "✘", ColorRed),
			Ui: &cli.BasicUi{
				Reader:      inStream,
				Writer:      outStream,
				ErrorWriter: errStream,
			},
		},
	}
}

func (u *ui) Status(prefix, message string, colorAttrs ColorAttrs) {
	u.withColoredPrefix(prefix, message, colorAttrs...)
}

func (u *ui) withColoredPrefix(prefix, msg string, attrs ...color.Attribute) {
	colored := color.New(attrs...).SprintfFunc()
	u.Output(fmt.Sprintf("%s  %s", colored("%12s", prefix), msg))
}

func uiPrefix(level, icon string, colorAttrs ColorAttrs) string {
	colored := color.New(colorAttrs...).SprintfFunc()
	return colored("%-6s%-2s  ", level, icon)
}
