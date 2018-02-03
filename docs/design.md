# Overview
Based on https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2 with the intent of
separating commands, library, and test code.

# Configuration
The `grs` program can be configured with a config file. The config file looks like
```$json
{
  "git": "/path/to/git_executable",
  "repos": [
    {"path":"/foo/bar/repo1"}
    {"path":"/home/repos/myproject"}
    ...
  ]
}
```

The config file can be read from
- `${GRS_CONF}`
- `$HOME/.grs.json`

The command line parameters for grs looks like
```$bash
Usage of grs:
  -repos rel/repo1:/abs/repo2:...
    target repos
  -verbose 
    verbose mode
```


# grs
Foundation for the "interesting" parts of the code. 
. logging
. abstraction layer for running shell commands
. implementations of commands to be eventually used inside by scripts

# script
The `grs` software consists of scripts that can be wired together through the go 
programming language. A `script` is a go function that runs with a context of 
"Command Runner" and "Repository Location". These two ideas are abstracted in order to test
scripts by mocking their dependencies.

# status
Contains data model for the UI - how statuses of multiple repositories are modeled inside the
program.

## User Relevant Statues

###ahead
The repo is unmodified and ahead of the remote repository. GRS will prompt user to push.

###behind
The repo is unmodified and behind the remote repository and cannot be cleanly fast-forwarded.
The user will be prompted to resolve conflicts.

###up-to-date
The repo is unmodified and neither ahead of behind the remote repository. Lets user know that
GRS has validated the repo recently.
 
###modified
The working directory of the repository contains uncomitted changes that can be cleanly applied
to the head of the remote repository.

###modified-conflict
The working directory of the repository contains uncommitted changes that will conflict with
the head of the remote repository.

###invalid
The status of the repository cannot be determined. Represents a generic error condition.


## Internal Statuses

RStat
 - dir: valid; invalid
 - branch: unknown; uptodate; ahead; behind; diverged;
 - index: unknown; unmodified;modified

 dir status | desc
-------------------
invalid     | The specified repo directory does not exist or is not a git repository
valid       | The specified repo exists and is a git repository

 branch/idx | unmodified | modified | unknown 
------------|------------|----------|---------
unknown     | notify     | "        | notify
up-to-date  | notify     | "        | notify
ahead       | notify     | "        | notify
behind      | rebase     | ???      | notify
diverged    | notify     | "        | notify

# config
The config module provides `GetCurrConfig`, which allows the user to specify the location of 
the `git` program and target repos using a config file.

# ctx
The ctx module is a singleton context that is available throughout the entire ilfe of the `grs`
program. The context (`GetContext()`) provides global methods:
```$golang
GetRepos() - Which repos to scan
GetGitExec() - The `git` executable to use
```
