package context

type RunnerContext interface {
	Version() string

	//Cwd() string

	ContextId() string

	Init(map[string]string)

	// LogDir: action write temp/log files to this dir
	// Suppose to be random generated dir in user's source tree
	LogDir() string

	/*
		   Actions need a lightweight way to pass results to downstream.

		   Complex data could be stored in files in LogDir and
			  pass filename as a result
	*/

	// Get one named value, first or last
	ActionResult(key string, asc bool) (value string)

	// Get all named value
	ActionResultAll(key string) (values []string)

	// Set current action
	SetCurrentResult(map[string]string) error
}
