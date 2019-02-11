package shexec

type AppContext struct {
	CommandRunner
	GitExec string
}

func newAppCtxWithDefaults() *AppContext {
	return &AppContext{
		CommandRunner: &ExecRunner{},
		GitExec:       "git",
	}
}

func NewAppContext(options ...AppContextOption) *AppContext {
	ctx := newAppCtxWithDefaults()
	for _, option := range options {
		option(ctx)
	}
	return ctx
}

// === START: options ===
type AppContextOption func(*AppContext)

func WithDefaultGitExec(gitExec string) AppContextOption {
	return func(ctx *AppContext) {
		ctx.GitExec = gitExec
	}
}

func WithCommandRunner(runner CommandRunner) AppContextOption {
	return func(ctx *AppContext) {
		ctx.CommandRunner = runner
	}
}

// === END: options ===
