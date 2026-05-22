# PowerShell argument completer for kxd. Mirrors scripts/_kxd_autocomplete.
#
# Setup: dot-source from your PowerShell profile ($PROFILE), after _kxd.ps1:
#   . "$HOME\bin\_kxd.ps1"
#   . "$HOME\bin\_kxd_autocomplete.ps1"
#
# Completes the first positional arg (a config filename) from `_kxd_prompt
# file list`. Calls the binary directly to avoid the wrapper's KUBECONFIG
# side effects during tab-completion.

Register-ArgumentCompleter -CommandName kxd -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    & _kxd_prompt.exe file list 2>$null |
        Where-Object { $_ -like "$wordToComplete*" } |
        ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
}
