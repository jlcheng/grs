package script

import "jcheng/grs/shexec"

type AppContext struct {
	shexec.CommandRunner
	GitExec string
}

func newAppCtxWithDefaults() *AppContext {
	return &AppContext{
		CommandRunner: &shexec.ExecRunner{},
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

func WithCommandRunner(runner shexec.CommandRunner) AppContextOption {
	return func(ctx *AppContext) {
		ctx.CommandRunner = runner
	}
}

// === END: options ===
