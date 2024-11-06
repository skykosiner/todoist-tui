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
	todoist := &Todoist{
		ctx:      ctx,
		config:   config,
		Projects: []Project{},
		Tasks:    []Task{},
	}

	todoist.getProjects()

	return todoist
}

func (t Todoist) makeRequest(path string) []byte {
	client := http.Client{}
	req, err := http.NewRequestWithContext(t.ctx, "GET", fmt.Sprintf("https://api.todoist.com/rest/v2/%s", path), nil)
	if err != nil {
		slog.Error("Error making request to todoist api", "error", err)
		return []byte{}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.config.Token))
	res, err := client.Do(req)
	if err != nil {
		if t.ctx.Err() == context.DeadlineExceeded {
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

func (t *Todoist) getProjects() {
	var projects []Project

	bytes := t.makeRequest("projects")
	if err := json.Unmarshal(bytes, &projects); err != nil {
		slog.Error("Erorr unmarashling JSON for projects", "error", err)
		return
	}

	t.Projects = projects
}

func (t *Todoist) GetTasksForProject(projectID string) {
	var tasks []Task

	bytes := t.makeRequest("tasks")
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		slog.Error("Erorr unmarashling JSON for tasks in project", "error", err, "project id", projectID)
		return
	}

	for _, task := range tasks {
		if task.ID != projectID {
		}
	}
}
