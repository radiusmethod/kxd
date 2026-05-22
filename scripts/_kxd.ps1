# PowerShell wrapper for kxd. Mirrors scripts/_kxd (the bash version).
#
# The Go binary (_kxd_prompt.exe) writes the user's selection to ~/.kxd. This
# wrapper then reads that file and sets $env:KUBECONFIG in the current
# PowerShell session — the same `source`-the-script trick the bash version
# uses, since a child process cannot mutate its parent's environment.
#
# Setup: dot-source this file from your PowerShell profile ($PROFILE):
#   . "$HOME\bin\_kxd.ps1"
# Then invoke as `kxd ...` like the POSIX version.

function kxd {
    & _kxd_prompt.exe @args

    $kxdFile = Join-Path $HOME '.kxd'
    if (-not (Test-Path $kxdFile)) {
        Remove-Item Env:KUBECONFIG -ErrorAction SilentlyContinue
        return
    }

    $selectedConfig = (Get-Content $kxdFile -Raw -ErrorAction SilentlyContinue)
    if ($null -ne $selectedConfig) { $selectedConfig = $selectedConfig.Trim() }

    if ([string]::IsNullOrEmpty($selectedConfig)) {
        Remove-Item Env:KUBECONFIG -ErrorAction SilentlyContinue
    } else {
        $env:KUBECONFIG = Join-Path $HOME ".kube\$selectedConfig"
    }
}
