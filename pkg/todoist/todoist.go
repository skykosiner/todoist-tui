package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	c "github.com/skykosiner/todoist-tui/pkg/config"
)

type DueDate struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
}

type Task struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"is_completed"`
	Due         DueDate
	ParentID    string `json:"parent_id"`
	Priority    string `json:"priority"`
	ProjectID   int    `json:"project_id"`
}

type Project struct {
	ID    string `json:"id"`
	Name  string `jsqn:"name"`
	Tasks []Task
}

type Todoist struct {
	ctx      context.Context
	config   c.Config
	Projects []Project
	Tasks    []Task
}

func NewTodoist(ctx context.Context, config c.Config) *Todoist {
	return &Todoist{
		ctx:      ctx,
		config:   config,
		Projects: getProjects(ctx, config),
		Tasks:    []Task{},
	}
}

func makeRequest(ctx context.Context, path string, token string) []byte {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.todoist.com/rest/v2/%s", path), nil)
	if err != nil {
		slog.Error("Error making request to todoist api", "error", err)
		return []byte{}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Request to Todoist API timed out")
		} else {
			slog.Error("Error making request to Todoist API", "error", err)
		}
		os.Exit(0)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		slog.Error("Error from Todoist API", "Status Code", res.StatusCode)
		os.Exit(0)
	}

	bytes, _ := io.ReadAll(res.Body)
	return bytes
}

func getProjects(ctx context.Context, config c.Config) []Project {
	var projects []Project

	bytes := makeRequest(ctx, "projects", config.Token)
	if err := json.Unmarshal(bytes, &projects); err != nil {
		slog.Error("Erorr unmarashling JSON for projects", "error", err)
		return projects
	}

	return projects
}
