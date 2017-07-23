package ui

import "github.com/fatih/color"

// ColorAttrs is a shell color attribute set
type ColorAttrs []color.Attribute

var (
	// ColorGray represents bold texts in the shell
	ColorGray = ColorAttrs{color.Bold}
	// ColorGreen represents green and bold texts in the shell
	ColorGreen = ColorAttrs{color.FgGreen, color.Bold}
	// ColorBlue represents blue and bold texts in the shell
	ColorBlue = ColorAttrs{color.FgBlue, color.Bold}
	// ColorYellow represents yellow and bold texts in the shell
	ColorYellow = ColorAttrs{color.FgYellow, color.Bold}
	// ColorRed represents red and bold texts in the shell
	ColorRed = ColorAttrs{color.FgRed, color.Bold}
)
