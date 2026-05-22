# kxd - Kubeconfig Switcher

<img src="assets/kxd.png" width="200">

kxd is a command-line utility that allows you to easily switch between Kubernetes configuration files (kubeconfig) contexts. This tool is designed to simplify the management of multiple Kubernetes clusters and contexts.

<img src="assets/demo.gif" width="500">

## Features

- Switch between different kubeconfig files.
- Switch between Kubernetes contexts within a kubeconfig file.
- Switch between Kubernetes namespaces in a context.

## Table of Contents

- [Installation](#installation)
    - [Homebrew](#homebrew)
    - [Makefile](#makefile)
    - [To Finish Installation](#to-finish-installation)
    - [Upgrading](#upgrading)
    - [Windows](#windows)
    - [Configuration](#configuration)
- [Usage](#usage)
    - [Switching Kubeconfig Files](#switching-kubeconfig-files)
    - [Switching Kubernetes Contexts](#switching-kubernetes-contexts)
    - [Switching Kubernetes Context Namespaces](#switching-kubernetes-context-namespaces)
    - [Getting Current Kubeconfig, Kubernetes Context or Context Namespace](#getting-current-kubeconfig-kubernetes-context-or-context-namespace)
    - [Version](#version)
    - [Persist KUBECONFIG across new shells](#persist-kubeconfig-across-new-shells)
    - [Show your set kubeconfig in your shell prompt](#show-your-set-kubeconfig-in-your-shell-prompt)
    - [Add autocompletion](#add-autocompletion)
    - [TL;DR (full config example)](#tldr-full-config-example)
- [Contributing](#contributing)
- [License](#license)

## Installation

Make sure you have Go installed. You can download it from [here](https://golang.org/dl/).

### Homebrew

```bash
brew tap radiusmethod/kxd
brew install kxd
```

or just

```bash
brew install radiusmethod/kxd/kxd
```
### Makefile

```bash
make install
```

### To Finish Installation
Add the following to your bash profile or zshrc then open new terminal or source that file

```sh
alias kxd="source _kxd"
```

Ex. `echo -ne '\nalias kxd="source _kxd"' >> ~/.zshrc`

### Upgrading
Upgrading consists of just doing a brew update and brew upgrade.

```sh
brew update && brew upgrade radiusmethod/kxd/kxd
```

### Windows

`kxd` is designed for POSIX shells (bash/zsh): the Go binary writes the user's selection to `~/.kxd`, then a wrapper script that you `source` reads that file and exports `KUBECONFIG` into your current shell. A child process can't mutate its parent's environment, so the wrapper indirection is mandatory — and that's what makes Windows non-trivial. Three paths work:

#### WSL (recommended)

From a WSL2 Ubuntu/Debian shell, follow the standard Linux instructions exactly: `make install`, then add `alias kxd="source _kxd"` to `~/.bashrc` or `~/.zshrc`. From WSL's perspective it's just Linux.

Caveat: the `KUBECONFIG` you set inside WSL is **not** visible to `kubectl.exe` invoked from PowerShell or `cmd`. Run `kubectl` from WSL too, or set the env var separately on the Windows side.

#### Git Bash / MSYS2

The Go binary cross-compiles cleanly and the bash wrapper is portable enough for Git Bash to `source`. Manual setup:

1. Build the Windows binary:
   ```sh
   GOOS=windows GOARCH=amd64 go build -o _kxd_prompt.exe .
   ```
2. Copy `_kxd_prompt.exe`, `scripts/_kxd`, and `scripts/_kxd_autocomplete` to a directory on your Git Bash `PATH` (e.g. `~/bin`).
3. Add to `~/.bashrc`:
   ```sh
   alias kxd="source _kxd"
   source _kxd_autocomplete
   ```
4. Make sure `~/.kube/` exists with your config files. In Git Bash, `~` resolves to `C:\Users\<you>`.

Untested by the maintainers. `~/.kube/config` symlinks created on the Windows side sometimes confuse path resolution.

#### Native PowerShell

`scripts/_kxd.ps1` and `scripts/_kxd_autocomplete.ps1` are PowerShell equivalents of the bash wrapper and autocomplete.

1. Build the Windows binary and put it somewhere on `$env:PATH`:
   ```powershell
   $env:GOOS = "windows"; $env:GOARCH = "amd64"
   go build -o _kxd_prompt.exe .
   # move _kxd_prompt.exe into e.g. C:\Users\<you>\bin
   ```
2. Copy `scripts/_kxd.ps1` and `scripts/_kxd_autocomplete.ps1` somewhere persistent (e.g. `C:\Users\<you>\bin`).
3. Dot-source both from your PowerShell profile (open it with `notepad $PROFILE`):
   ```powershell
   . "$HOME\bin\_kxd.ps1"
   . "$HOME\bin\_kxd_autocomplete.ps1"
   ```
4. Restart PowerShell. `kxd` is now a function in your session.

The function reads/writes `$HOME\.kxd` and `$HOME\.kube\<name>` — same layout as the POSIX version, so configs interoperate with WSL or Git Bash on the same machine if you point them at the same `.kube` directory.

## Configuration

By default, Kubeconfig Switcher looks for files with an extension of `.conf`. You can customize the behavior by setting an environment variable.
This can be a single matcher or a comma seperated string for multiple matchers.

- `KXD_MATCHER`: The file matcher(s) used to identify kubeconfig files (default is `.conf`).

## Usage

 * See docs for more info [kxd](docs/kxd.md)

### Switching Kubeconfig Files

It is possible to shortcut the menu selection by passing the config name you want to switch to as an argument.

```bash
> kxd dev.conf
Config dev.conf set.
```

To switch between different kubeconfig files using the menu, use the following command:

```bash
kxd f s
```

This command will display a list of available kubeconfig files in your `~/.kube` directory. Select the one you want to use.

### Switching Kubernetes Contexts

To switch between Kubernetes contexts within a kubeconfig file, use the following command:

```bash
kxd ctx s
```

This command will display a list of available contexts in your current kubeconfig file. Select the one you want to switch to.

### Switching Kubernetes Context Namespaces

To switch between Kubernetes context namespaces within a kubeconfig context, use the following command:

```bash
kxd ns s
```

This command will display a list of kubernetes namespaces in your currently set cluster. Select the one you want to switch to.


### Getting Current Kubeconfig, Kubernetes Context or Context Namespace

To get the currently set Kubeconfig, Kubernetes Context or Context Namespace, use the following commands:

```bash
kxd f c
```

This command will display the currently set kubeconfig file.

```bash
kxd ctx c
```

This command will display the currently set Kubernetes Context.

```bash
kxd ns c
```

This command will display the currently set Kubernetes Context Namespace.

### Version

To check the version of Kubeconfig Switcher, use the following command:

```bash
kxd version
```

## Persist KUBECONFIG across new shells
To persist the set config when you open new terminal windows, you can add the following to your bash profile or zshrc.

```bash
export KUBECONFIG=$(kxd file current)
```

### Show your set kubeconfig in your shell prompt
For better visibility into what your shell is set to it can be helpful to configure your prompt to show the value of the env variable `KUBECONFIG`.

<img src="assets/screenshot.png" width="700">

Here's a sample of my zsh prompt config using oh-my-zsh themes

```sh
# Kubeconfig info
local kxd_info='$(kxd_config)'
function kxd_config {
  local config="${KUBECONFIG:=}"
    if [ -z "$config" ]
    then
          echo -n ""
    else
          config=$(basename $config)
          echo -n "%{$fg_bold[blue]%}kx:(%{$fg[cyan]%}${config}%{$fg_bold[blue]%})%{$reset_color%} "
    fi
}
```

```sh
PROMPT='OTHER_PROMPT_STUFF $(kxd_info)'
```

To include prompt support in OhMyZsh, add the following lines to your `~/.p10k.zsh` file:

```sh
# kxd prompts
typeset -g _kxd_config
typeset -g _kxd_basename=''
typeset -g _kxd_content="kx:(${_kxd_basename})"

function prompt_kxd() {
local _kxd_config="${KUBECONFIG:=}"
if [ -z "$_kxd_config" ]
then
  _kxd_basename=''
else
  _kxd_basename="%F{cyan}$(basename $_kxd_config)%f"
fi

_kxd_content="kx:(${_kxd_basename})"
p10k segment -b 0 -f 4 -t ${_kxd_content}
}

function instant_prompt_kxd() {
p10k segment -b 0 -f 4 -t ${_kxd_content}
}
```

Then add `kxd` to either your left or right prompt segments.

<img src="assets/ohmyzsh-screenshot.png" width="700">

## Add autocompletion
You can add autocompletion when passing config as argument by adding the following to your bash profile or zshrc file.
`source _kxd_autocomplete`

Now you can do `kxd my-k` and hit tab and if you had a config `my-kubeconfig` it would autocomplete and find it.

## TL;DR (full config example)
```sh
alias kxd="source _kxd"
source _kxd_autocomplete
export KXD_MATCHER="-config,.conf"
export KUBECONFIG=$(kxd file current)
```

## Contributing

If you encounter any issues or have suggestions for improvements, please open an issue or create a pull request on [GitHub](https://github.com/radiusmethod/kxd).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
