## v0.2.0 (August 5, 2026)

**Breaking.** kxd is now a single binary named `kxd`, set up with one line in your shell rc file. Pre-1.0, so no compatibility shims.

* Added `kxd init <shell>` — `eval "$(kxd init zsh)"` now replaces the `kxd` alias, the completion `source`, and the `export KUBECONFIG=$(kxd file current)` persistence line. Supports bash, zsh, fish, and PowerShell.
* Added `kxd shellenv [shell]`, which prints the export/unset statement for the active kubeconfig. This is what the generated function evals.
* Added fish support. PowerShell setup is now one line instead of building a binary and dot-sourcing two scripts by hand.
* The binary is installed as `kxd`, not `_kxd_prompt`. The generated integration calls it through `command kxd`, which skips the shell function of the same name.
* Removed `scripts/` entirely: `_kxd`, `_kxd_autocomplete`, `_kxd.ps1`, and `_kxd_autocomplete.ps1`. `make install` now installs one file.
* `~/.kxd` parsing and the `~/.kube` path join now live only in the Go code, instead of being reimplemented in every wrapper script.
* Values are shell-quoted, so kubeconfig filenames containing spaces or quotes work.
* Completion now covers the `file`/`context`/`namespace` subcommands and their arguments, not just top-level config names. Namespace values are excluded on purpose, since listing them calls the live cluster.
* **Behavior change:** `kxd <unknown-config>` now writes its warning to stderr and exits 1, instead of writing to stdout and exiting 0. Scripts relying on the old exit code need updating. This keeps ANSI color codes out of the command substitutions the shell integration evals.

Releases are now built by GoReleaser:

* Releases ship prebuilt binaries for macOS, Linux, and Windows on amd64 and arm64, with `checksums.txt`. Installing no longer compiles from source, and Windows binaries are published for the first time.
* Homebrew distribution moves from a formula to a **cask**, generated on each tag. Reinstall with `brew uninstall kxd && brew install radiusmethod/kxd/kxd` if brew complains about the change.
* `kxd version` now reports the git tag, injected at build time. `make install` derives it from `git describe`, and a plain `go build` reports `dev`.
* Prerelease tags must now be SemVer-hyphenated (`v0.3.0-beta1`), not `v0.3.0beta`.

To upgrade: replace your shell config lines as described in "Upgrading from pre-v0.2.0" in the README, then delete any leftover `_kxd_prompt`, `_kxd`, and `_kxd_autocomplete` files. Your `~/.kxd` file is unchanged and carries over.

## v0.1.4 (May 22, 2026)
* Updated Kubernetes client libraries (client-go v0.36.1), cobra (v1.10.2), and Go toolchain to 1.26.

## v0.1.3 (August, 10, 2024)
* Fix issue where wrong index is selected.

## v0.1.2 (August, 10, 2024)
* Fixed a crash when scrolling through an empty list, which occurred when searches returned no results.

## v0.1.1 (August, 1, 2024)
* Adds circular scrolling #39.

## v0.1.0 (January 26, 2024)
* Added namespace command.

## v0.0.9 (December 23, 2023)
* Added autocomplete script to install.

## v0.0.8 (October 25, 2023)
* Fix for setting default as argument.
* Changed the way kxd file current works to check `~/.kxd` file then default to `~/.kube/config`.

## v0.0.7 (October 23, 2023)
* Added support for setting config names as argument.
* Added list command to `kxd file` and `kxd ctx`.
* The root command now defaults to `kxd file switch` if no sub-commands are passed in.

## v0.0.6 (September 29, 2023)
* Allow for specifying multiple matchers as a comma seperated string for `KXD_MATCHER` environment variable.

## v0.0.5 (September 25, 2023)
* Small fix for running `kxd file switch -h`.

## v0.0.4 (September 25, 2023)
* Allow for listing of current config and context.

## v0.0.3 (September 23, 2023)
* Added support for default config at `~/.kube/config`.
* Added context switcher.

## v0.0.2 (September 22, 2023)
* Added `KXD_MATCHER` environment variable for specifying different config matchers.

## v0.0.1 (September 21, 2023)
* Initial Release
