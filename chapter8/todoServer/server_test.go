package main

import (
	"bytes"
	"clis-in-go/chapter2/todo"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	// Create a temporary todofile for use in the testing
	tempTodoFile, err := os.CreateTemp("", "todotest")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(newMux(tempTodoFile.Name()))

	// Add a few items for testing
	for i := 1; i < 3; i++ {
		var body bytes.Buffer

		taskName := fmt.Sprintf("Task number %d.", i)
		item := struct {
			Task string `json:"task"` // the JSON tag sets the key in the returned object to lowercase t
		}{
			Task: taskName,
		}

		// the equivalent of body = json.loads(item); in python
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}

	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
	}
}

func TestComplete(t *testing.T) {
	url, cleanup := setupAPI(t)
	defer cleanup()

	t.Run("Complete", func(t *testing.T) {
		u := fmt.Sprintf("%s/todo/1?complete", url)

		req, err := http.NewRequest(http.MethodPatch, u, nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected %q, got %q", http.StatusText(http.StatusNoContent), http.StatusText(resp.StatusCode))
		}
	})

	t.Run("CheckComplete", func(t *testing.T) {
		r, err := http.Get(url + "/todo")
		if err != nil {
			t.Error(err)
		}

		if r.StatusCode != http.StatusOK {
			t.Fatalf("Expected %q, got %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
		}

		var resp todoResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if len(resp.Results) != 2 {
			t.Errorf("Expected 2 items, got %d", len(resp.Results))
		}
		if !resp.Results[0].Done {
			t.Errorf("Expected %t, but got %t", true, resp.Results[0].Done)
		}
		if resp.Results[1].Done {
			t.Errorf("Expected %t, but got %t", false, resp.Results[1].Done)
		}

	})

}

func TestDelete(t *testing.T) {
	url, cleanup := setupAPI(t)
	defer cleanup()

	t.Run("Delete", func(t *testing.T) {
		u := fmt.Sprintf("%s/todo/1", url)
		req, err := http.NewRequest(http.MethodDelete, u, nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected %q, got %q", http.StatusText(http.StatusNoContent), http.StatusText(resp.StatusCode))
		}
	})

	t.Run("CheckDelete", func(t *testing.T) {
		r, err := http.Get(url + "/todo")
		if err != nil {
			t.Error(err)
		}

		if r.StatusCode != http.StatusOK {
			t.Fatalf("Expected %q, got %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
		}

		var resp todoResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if len(resp.Results) != 1 {
			t.Errorf("Expected 1 item, got %d", len(resp.Results))
		}

		expTask := "Task number 2."
		if resp.Results[0].Task != expTask {
			t.Errorf("Expected %q, but got %q", expTask, resp.Results[0].Task)
		}
	})
}

func TestAdd(t *testing.T) {
	url, cleanup := setupAPI(t)
	defer cleanup()

	taskName := "Task number 3."
	t.Run("Add", func(t *testing.T) {
		var body bytes.Buffer
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		// Load the value of body into the item struct
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		// Make the request
		r, err := http.Post(url+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the status code is as expected
		if r.StatusCode != http.StatusCreated {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusCreated), http.StatusText(r.StatusCode))
		}
	})

	t.Run("CheckAdd", func(t *testing.T) {
		r, err := http.Get(url + "/todo/3")
		if err != nil {
			t.Fatal(err)
		}

		// Verify the status code is as expected
		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
		}

		var resp todoResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if resp.Results[0].Task != taskName {
			t.Errorf("Expected %q, got %q.", taskName, resp.Results[0].Task)
		}
	})
}

func TestGet(t *testing.T) {
	var (
		resp struct {
			Results      todo.List `json:"results"`
			Date         int64     `json:"date"`
			TotalResults int       `json:"total_results"`
		}
		// body []byte
		// err  error
	)

	testCases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{
			name:       "GetRoot",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There's an API here",
		},
		{
			name:    "NotFound",
			path:    "/todo/500",
			expCode: http.StatusNotFound,
		},
		{
			name:       "GetAll",
			path:       "/todo",
			expCode:    http.StatusOK,
			expItems:   2,
			expContent: "Task number 1.",
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			r, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close()

			if r.StatusCode != tc.expCode {
				t.Fatalf("Expected %q, got %q.", http.StatusText(tc.expCode),
					http.StatusText(r.StatusCode))
			}

			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}

				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expected %q, got %q.", tc.expContent,
						string(body))
				}
			case strings.Contains(r.Header.Get("Content-Type"), "application/json"):
				if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Error(err)
				}

				if resp.TotalResults != tc.expItems {
					t.Errorf("Expected %d items, got %d.", tc.expItems, resp.TotalResults)
				}
				if resp.Results[0].Task != tc.expContent {
					t.Errorf("Expected %q, got %q.", tc.expContent, resp.Results[0].Task)
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}
		})
	}
}
