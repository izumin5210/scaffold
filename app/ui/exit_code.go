package ui

import (
	"github.com/mitchellh/cli"
)

// Exit codes are values representing an exit code for a error type.
const (
	ExitCodeOK int = 0

	ExitCodeUnknownError = 10 + iota
	ExitCodeScaffoldNameMissingError
	ExitCodeScffoldNotFoundError
	ExitCodeFailedToGetScaffoldsError
	ExitCodeFailedToCreatetScaffoldsError

	ExitCodeInvalidArgumentListLengthError = cli.RunResultHelp
)
