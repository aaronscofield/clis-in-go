# Word Counter

## Description
A basic word counter CLI utility.

## Flags
- `-l` - Count lines of input from stdin (default False)
- `-b` - Count bytes of input from stdin (default False)

## Usage Examples
Counting words:
`echo "My first command works?" | ./wc`

Counting lines:
`cat main.go | ./wc -l`

Counting bytes:
`cat main.go | ./wc -b`

Building for windows:
```
export GOOS=windows go build
go build
```

## Run tests
`go test -v`