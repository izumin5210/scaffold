package app

import (
	"io"

	"fmt"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
)

// UIColor is a shell color attribute set
type UIColor []color.Attribute

var (
	// UIColorGray represents bold texts in the shell
	UIColorGray = UIColor{color.Bold}
	// UIColorGreen represents green and bold texts in the shell
	UIColorGreen = UIColor{color.FgGreen, color.Bold}
	// UIColorBlue represents blue and bold texts in the shell
	UIColorBlue = UIColor{color.FgBlue, color.Bold}
	// UIColorYellow represents yellow and bold texts in the shell
	UIColorYellow = UIColor{color.FgYellow, color.Bold}
	// UIColorRed represents red and bold texts in the shell
	UIColorRed = UIColor{color.FgRed, color.Bold}
)

// UI is an interface for intaracting with the shell
type UI interface {
	cli.Ui
	Status(prefix, message string, uiColor UIColor)
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

func (u *ui) Status(prefix, message string, uiColor UIColor) {
	u.withColoredPrefix(prefix, message, uiColor...)
}

func (u *ui) withColoredPrefix(prefix, msg string, attrs ...color.Attribute) {
	colored := color.New(attrs...).SprintfFunc()
	u.Output(fmt.Sprintf("%s  %s", colored("%12s", prefix), msg))
}
