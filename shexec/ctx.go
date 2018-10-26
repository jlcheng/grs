package shexec

type AppContext struct {
	CommandRunner
	defaultGitExec  string
}

func NewAppContext() *AppContext {
	return &AppContext{
		defaultGitExec:  "git",
	}
}
func NewAppContextWithRunner(runner CommandRunner) *AppContext {
	return &AppContext{
		CommandRunner:   runner,
		defaultGitExec:  "git",
	}
}

func (ctx *AppContext) GetGitExec() string {
	return ctx.defaultGitExec
}

func (ctx *AppContext) SetGitExec(defaultGitExec string) {
	ctx.defaultGitExec = defaultGitExec
}
