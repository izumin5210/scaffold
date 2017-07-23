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
	Status(prefix, message string, colorAttrs ColorAttrs)
}

type ui struct {
	cli.Ui
}

// NewUI creates a new UI instance from streams
func NewUI(inStream io.Reader, outStream, errStream io.Writer) UI {
	return &ui{
		Ui: &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
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
