function build {
    echo building ./cmd/$1
    compile $1
}


function run {
    echo running ./cmd/$1 ${@:2}
    $(compile $@)
}

function compile {
    mkdir -p ./target/$1
    go build -o ./target/$1/avisha ./cmd/$1
    echo target/$1/avisha ${@:2}
}

function watch {
    # https://github.com/watchexec/watchexec
    watchexec -i target "go build -o ./target/$1/avisha ./cmd/$1" &
    watchexec  -w target -i target/*.json -r "./target/$1/avisha"
}