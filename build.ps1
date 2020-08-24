function build {
    param (
        [Parameter(Mandatory = $true)]
        [string]
        $bin
    )
    compile $bin
}

# Compile and execute.
function run {
    param (
        [Parameter(Mandatory = $true)]
        [string]
        $bin
    )
    Invoke-Expression $(compile $bin)
}

# Compile the specified binary. 
# Binary must be a main package under the "cmd" directory.
function compile {
    param (
        [Parameter(Mandatory = $true)]
        [string]
        $bin
    )
    mkdir -Force -p ./target/$bin | Out-Null
    go build -o ./target/$bin/avisha.exe ./cmd/$bin
    return "./target/$bin/avisha.exe"
}

function watch {
    param (
        [Parameter(Mandatory = $true)]
        [string]
        $bin
    )
    watchexec.exe -i target/* "go build -o ./target/$bin/avisha.exe ./cmd/$bin" &
    watchexec.exe -w target -i target/*.json -r ".\target\$bin\avisha.exe"
}
