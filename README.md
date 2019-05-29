# GRS (Grease)

Grs is a dashboard for your Git repos and keep your repos up-to-date. You can think of it as a glorified cron job with a
few differences.

1. Grs will try to auto-rebase your local changes. If there is a conflict, Grs will keep your local files modified.
2. Grs can optionally auto-push your local changes.
3. Grs provides a dashboard view of your Git repos.

I use Grs as a replacement for Dropbox and Evernote. Many files that matter to me are text files: scripts, org-mode
documents, and journals. Thus, it is natural for me to use Git to manage them. Git does a good job of auto-resolving
simple conflicts where possible. As an engineer, this tool gives me control. I can pay for Git hosting and take
ownership of my data.

I usually run Grs in a terminal window and just let it do its thing:
```
  Grs [May 29 10:28:12AM PST]──────────────────────────────────────────────────────────────────────
  repo [/home/jcheng/repos/repo1]⯅ status is BRANCH_UPTODATE, INDEX_UNMODIFIED, 20 hours ago.
  repo [/home/jcheng/repos/repo2] status is BRANCH_AHEAD, INDEX_MODIFIED, 77 seconds ago.
  repo [/home/jcheng/workspace/repo3]⯅ status is BRANCH_UPTODATE, INDEX_UNMODIFIED, 4 days ago.
  repo [/home/jcheng/repo4]⯅ status is BRANCH_UPTODATE, INDEX_UNMODIFIED, 55 minutes ago.
...
```

# tldr;
Go version 1.12+ is required. Grs uses Go modules for dependencies management and 1.12-specific APIs.

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
~/bin/grs
```

To manually refresh repo status, simply hit CTRL-R

