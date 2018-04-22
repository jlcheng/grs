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
unknown     | Likely because the repo directory is invalid, rarely used
uptodate    | Local branch is in sync with remote branch
ahead       | Local branch has changes 
behind      | Remote branch has changes
diverged    | Both sides are _known_ to have changes
untracked   | Local branch does not have a remote branch


# config
The config module provides access to user config files. Meant to be encapsulated by the AppContext 
facade.

# ctx and sctx
sctx - The Script Context describes dependencies between top-level script components: inputs
from the command line and the AppContext object.

ctx - The AppContext instance describes dependencies at a level lower than Script Context.
It provides globally useful methods:

```$golang
GetGitExec() - The `git` executable to use
GetPrinter() - (TODO) The API for presenting messages to users with different verbosity settings
DB() - Poorly named. Holds metadata on the user's repos, e.g., "last attempt to auto-merge/rebase."
DBService() - The API to persist the DB object to disk (using a Key-Value paradigm). Should be renamed to DBDAO, which 
implies the "DB" abstraction should be renamed.
ConfParams() - Parameters loaded from disk; Acts as a facade for loading anything considered to be user perference or
 configuration.
```
