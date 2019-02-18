# GRS (Grease)

Gfrs is a command-line program that polls your Git repos and keep them up to date. You can think of it as a glorified
cron job with a few differences.

1. Grs will try to auto-rebase your local changes. If there is a conflict, Grs will keep your local files unchanged.
2. Grs can optionally auto-push your changes.
3. Grs provides a dashboard view of your Git repos.

I use Grs as a replacement for Dropbox and Evernote. Many files that matter to me are text files: scripts, org-mode
documents, and journals. Thus, it is natural for me to use Git to manage them. Git does a good job of auto-resolving
simple conflicts where possible. As an engineer, this tool gives me control. I can pay for Git hosting and take
ownership of my data.

I usually run Grs in a terminal window and just let it do its thing:
```
=== Feb 13 9:02PM PST ===
repo [/home/jcheng/privprjs/dotfiles] status IS UP-TO-DATE, UNMODIFIED, 64 minutes ago.
repo [/home/jcheng/privprjs/forget] status IS UP-TO-DATE, UNMODIFIED, 4 days ago.
repo [/home/jcheng/privprjs/grs] status IS UP-TO-DATE, MODIFIED, 13 minutes ago.
repo [/home/jcheng/privprjs/playground] status IS UP-TO-DATE, UNMODIFIED, 2 days ago.
repo [/home/jcheng/org] status IS UP-TO-DATE, UNMODIFIED, 2 seconds ago.
```

# tldr;
Go version 1.11+ is required. Grs uses Go modules for dependencies management.

Install Grs
```
$ make all     # Runs tests and creates out/grs
$ make install # Installs grs in $HOME/bin
```

Create a configuration file in ~/.grs.toml
```
cat<<ENDL
# Tells Grs to run git pull in these three directories
repos = [
  "/home/jcheng/repos/grs",
  "/home/jcheng/repos/git",
  "/home/jcheng/repos/foo",
]

# Tells Grs it may perform git push using auto-generated commit messages
[[repo_config]]
id = "/home/jcheng/repos/foo"
push_allowed = true
ENDL
```

Finally, run Grs in a terminal window
```
~/bin/grs --use-cui
```
