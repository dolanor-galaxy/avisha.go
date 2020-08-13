function build {
    echo building ./cmd/$1
    compile $1
}


function run {
    echo running ./cmd/$1
    $(compile $1)
}

function compile {
    mkdir -p ./target/$1
    go build -o ./target/$1/avisha ./cmd/$1
    echo target/$1/avisha
}