# Purpose
The `grs` program:
- Generates a dashboard of your most commonly used repositories and their status.
- Automatically incorporates latest changes from upstream repository.

I hope it will be more useful than running `git fetch` and `git status` in cron. 

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

The `grs` configuration defaults to `$HOME/.grs.json` but can be overriden with the `${GRS_CONF}` variable.

# Usage
The command line parameters for grs looks like
```$bash
Usage of grs:
  -verbose 
    verbose mode
```

