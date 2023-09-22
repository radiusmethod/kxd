# kxd

Kubeconfig Switcher in Go

Easily switch between Kubeconfigs

<img src="assets/demo.gif" width="500">

## Requirements
min go 1.16

## Install

### Homebrew

```sh
brew tap radiusmethod/kxd
brew install kxd
```

or just

```sh
brew install radiusmethod/kxd/kxd
```
### Makefile

```sh
make install
```



### To Finish Installation
Add the following to your bash profile or zshrc then open new terminal or source that file

```sh
alias kxd="source _kxd"
```

Ex. `echo -ne '\nalias kxd="source _kxd"' >> ~/.zshrc`

## Usage
```sh
kxd
```

## Show your set kubeconfig in your shell prompt
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

Looking for a naming convention besides `.conf`. You can set an environment variable that lets you specify a string 
matcher. For example, lets say your files all end in `-config`. You'd set `KXD_MATCHER="-config"`
