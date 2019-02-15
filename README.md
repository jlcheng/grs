# GRS (Grease)

Grease is a cron-job like daemon that polls your git repos and keep them up to date. You can think of it as a glorified
cron job with a few differences.

1. Grease will try to auto-rebase your local changes. If a conflict arises, Grease will keep your local files unchanged.
2. Grease can optionally auto-push your changes.

I use it as a replacement for Dropbox and Evernote. Many files that matter to me are text files: scripts, org-mode
documents, journals. Git does a good job of auto-resolving simple conflicts where possible. It helps that I am familiar
with Git, which allows me to effectively resolve conflicts. Finally, as an engineer, this tool gives me more control
over where my data will be stored. I can pay for git hosting and take ownership over my data.

I usually run it in a terminal window and just let it do its thing:
```
=== Feb 13 9:02PM PST ===
repo [/home/jcheng/privprjs/dotfiles] status IS UP-TO-DATE, UNMODIFIED, 64 minutes ago.
repo [/home/jcheng/privprjs/forget] status IS UP-TO-DATE, UNMODIFIED, 4 days ago.
repo [/home/jcheng/privprjs/grs] status IS UP-TO-DATE, MODIFIED, 13 minutes ago.
repo [/home/jcheng/privprjs/playground] status IS UP-TO-DATE, UNMODIFIED, 2 days ago.
repo [/home/jcheng/org] status IS UP-TO-DATE, UNMODIFIED, 2 seconds ago.
```

# tldr;
Install Grease
```
$ make all     # Runs tests and creates out/grs
$ make install # Installs grs in $HOME/bin
```

Create a configuration file in ~/.grs.toml
```
cat<<ENDL
# Tells Grease to run git pull in these three directories
repos = [
  "/home/jcheng/repos/grs",
  "/home/jcheng/repos/git",
  "/home/jcheng/repos/foo",
]

# Tells Grease it may perform git push using auto-generated commit messages
[[repo_config]]
id = "/home/jcheng/repos/foo"
push_allowed = true
ENDL
```