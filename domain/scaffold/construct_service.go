package scaffold

// ConstructService creates directories and files from scaffold templates
type ConstructService interface {
	Perform(
		sc Scaffold,
		v interface{},
		cb ConstructCallback,
		cCb ConstructConflictedCallback,
	) error
}
