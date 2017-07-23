package app

import (
	"io"

	"fmt"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
)

type UIColor []color.Attribute

var (
	UIColorGray  = UIColor{color.Bold}
	UIColorGreen = UIColor{color.FgGreen, color.Bold}
	UIColorBlue  = UIColor{color.FgBlue, color.Bold}
	UIColorYello = UIColor{color.FgYellow, color.Bold}
	UIColorRed   = UIColor{color.FgRed, color.Bold}
)

type UI interface {
	cli.Ui
	Status(prefix, message string, uiColor UIColor)
}

type ui struct {
	cli.Ui
}

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

func (u *ui) Status(prefix, message string, uiColor UIColor) {
	u.withColoredPrefix(prefix, message, uiColor...)
}

func (u *ui) withColoredPrefix(prefix, msg string, attrs ...color.Attribute) {
	colored := color.New(attrs...).SprintfFunc()
	u.Output(fmt.Sprintf("%s  %s", colored("%12s", prefix), msg))
}
