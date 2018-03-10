# Overview
Based on https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2 with the intent of
separating commands, library, and test code.

# grs
Foundation for the "interesting" parts of the code. 
. logging
. abstraction layer for running shell commands
. implementations of commands to be eventually used inside by scripts

# script
The `grs` software consists of scripts that can be wired together. A `script` is a function
that accepts "application context", "command runner", and "repository location". 

# status
Contains data model for the UI: How statuses of multiple repositories are modeled inside the
program.

`RStat`
 - dir: valid; invalid
 - branch: unknown; uptodate; ahead; behind; diverged; untracked
 - index: unknown; unmodified;modified

 dir status | desc
-------------------
invalid     | The specified repo directory does not exist or is not a git repository
valid       | The specified repo exists and is a git repository

 branch status | desc
-------------------
unknown     | Likely because the repo directory is invalid
uptodate    | Local branch is in sync with remote branch
ahead       | Local branch has changes 
behind      | Remote branch has changes
diverged    | Both sides are _known_ to have changes
untracked   | Local branch does not have a remote branch


# config
The config module provides `GetCurrConfig`, which allows the user to specify the location of 
the `git` program and target repos using a config file.

# ctx
The AppContext instance is available throughout most of the `grs` program. It provides globally
useful methods:

```$golang
GetRepos() - Which repos to scan
GetGitExec() - The `git` executable to use
GetPrinter() - (TODO) The API for presenting messages to users with different verbosity settings
GetDB() - Read/Save the "last fetched" time for each Repo
```
