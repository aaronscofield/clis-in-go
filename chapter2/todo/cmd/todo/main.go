package main

import (
	"bufio"
	"clis-in-go/chapter8/todo"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const FILE_ENV_VAR = "TODO_FILENAME"

func main() {
	// Use the flag package to facilitate command args
	// This also enables a "-h" arg by default, listing the below defined flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item index to be completed")
	delete := flag.Int("delete", 0, "Item index to be deleted")
	verbose := flag.Bool("verbose", false, "Show verbose output for tasks")
	showCompleted := flag.Bool("completed", false, "Show completed tasks in the list")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2025\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check environment variable for file path
	var todoFileName = ".todo.json"
	value, found := os.LookupEnv(FILE_ENV_VAR)
	if found {
		todoFileName = value
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		// By specifying exit code 1, it prints the value to STDERR
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list of items
	case *list:
		fmt.Print(l.String(*verbose, *showCompleted))
	case *add:
		// Determine whether to add the task via args or STDIN
		t, err := getTask(os.Stdin, flag.Args()...)

		// If an error occured, print it and exit
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Add the task to the list

		l.Add(t)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *complete > 0:
		// Mark the task as complete
		l.Complete(*complete)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		// Delete the task from the list
		l.Delete(*delete)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the description for a new task from:
// arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
