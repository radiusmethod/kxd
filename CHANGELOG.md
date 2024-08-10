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
