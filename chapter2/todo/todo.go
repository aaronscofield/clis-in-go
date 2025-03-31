package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

// Interface (Not required)
// type Stringer interface {
// 	 String() string
// }

// String prints out a formatted list of tasks
func (l *List) String(verbose bool, showCompleted bool) string {
	formatted := ""

	for k, t := range *l {
		prefix := "  "
		if t.Done {
			if showCompleted {
				prefix = "X "
			} else {
				continue
			}
		}

		formatted += fmt.Sprintf("%s%d: %s ", prefix, k+1, t.Task)

		if verbose {
			formatted += fmt.Sprintf("(%s) (%s)", t.CreatedAt, t.CompletedAt)
		}
		formatted += "\n"

	}

	return formatted
}

// Since item starts wiht a lowercase "i", it is private to the package.
// If it had started with an uppercase "I", it is considered public
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

// By using (l *List), it says that Add is called on a list, like so:
// ex. myList.Add("task");
// This is called a "reciever"
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist in the list", i)
	}

	ls[i-1].CompletedAt = time.Now()
	ls[i-1].Done = true

	return nil
}

func (l *List) Delete(i int) error {
	ls := *l

	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist in the list", i)
	}

	// Using the slices import
	*l = slices.Delete(ls, i-1, i)

	// Using append
	// *l = append(ls[:i-1], ls[i:]...)

	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}
