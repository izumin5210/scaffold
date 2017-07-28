package ui

import (
	"github.com/mitchellh/cli"
)

// Exit codes are values representing an exit code for a error type.
const (
	ExitCodeOK int = 0

	ExitCodeUnknownError = 10 + iota
	ExitCodeScaffoldNameMissingError
	ExitCodeFailedToGetScaffoldsError
	ExitCodeFailedToCreatetScaffoldsError

	ExitCodeInvalidArgumentListLengthError = cli.RunResultHelp
	ExitCodeScffoldNotFoundError           = cli.RunResultHelp
)
