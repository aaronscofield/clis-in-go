package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

/*
 * Integration tests to compile and execute commands of the CLI
 */

const FILE_ENV_VAR = "TODO_FILENAME"

var (
	binName  = "todo"
	fileName = ".todo.json"
)

// Setup function to build the tool and teardown
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	// Add .exe extension for Windows OS
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	value, found := os.LookupEnv(FILE_ENV_VAR)
	if found {
		fileName = value
	}

	// Build the utility and verify that it built successfully
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	task := "test task number 1"

	// Test the creation of a new task
	t.Run("AddNewTaskFromArgs", func(t *testing.T) {
		// Add task through the args
		cmd := exec.Command(cmdPath, "-add", task)

		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		list := exec.Command(cmdPath, "-list")
		out, err := list.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s \n", task)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromStdin", func(t *testing.T) {
		// Add task through the standard input
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()

		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		list := exec.Command(cmdPath, "-list")
		out, err := list.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s \n  2: %s \n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	// Test the retrieval of a task list
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s \n  2: %s \n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {

		// Mark the first task in the list as complete
		cmd := exec.Command(cmdPath, "-complete", "1")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// Retrieve the list of tasks
		listcmd := exec.Command(cmdPath, "-list")
		out, err := listcmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// Expect the list of tasks to have two tasks, the first being completed
		expected := fmt.Sprintf("  2: %s \n", task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {

		// Delete the first task in the slice
		cmd := exec.Command(cmdPath, "-delete", "1")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// Retrieve the list of tasks
		listcmd := exec.Command(cmdPath, "-list")
		out, err := listcmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// Expect the list of tasks to have only the second task
		expected := fmt.Sprintf("  1: %s \n", task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
