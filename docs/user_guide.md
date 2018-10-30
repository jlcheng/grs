# Purpose
The `grs` program:
- Generates a dashboard of your most commonly used repositories and their status.
- Automatically incorporates latest changes from upstream repository.

I hope it will be more useful than running `git fetch` and `git status` in cron. 

# Configuration
The `grs` program can be configured with a config file. The config file looks like
```$toml
# $HOME/.grs.toml
repos = [
  "/home/jcheng/org",
  "/home/jcheng/dotfiles",
  "/home/jcheng/code",
]

[[repo_config]] # allow grs to push to the org repo
id = "/home/jcheng/org"
push_allowed = true

# The default configuration is push_allowed = false
# [[repo_config]]
# id = "/home/jcheng/code"
# push_allowed = false
```
