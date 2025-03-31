# Todo Task List

## Description
A CLI for interactions with a basic task list

## Flags
- `-add` - Add a task to the list via STDIN or args
- `-list` - List existing tasks in the list
- `-verbose` - Include information about creation/completion datetime in the list output
- `-completed` - Include completed tasks in the list output
- `-complete` - Mark a task as complete by index
- `-delete` - Delete a task from the list

## Usage Examples
List tasks, including completed ones, with verbose output:
`./todo --list -completed -verbose`

Add a task using the CLI arguments:
`./todo --add "Secondtask"`

Add a task using STDIN:
`echo "From stdin" | ./todo --add`

Delete a task by index:
`./todo --delete 1`

Mark a task as completed:
`./todo --complete 1`

Building for windows:
```
export GOOS=windows go build
go build
```

## Run tests
`go test -v`