package main

import (
	"context"
	"log/slog"
	"time"

	c "github.com/skykosiner/todoist-tui/pkg/config"
	t "github.com/skykosiner/todoist-tui/pkg/todoist"
)

func main() {
	config, err := c.NewConfig()
	if err != nil {
		slog.Error("Error getting config", "error", err)
		return
	}

	// Give the requests each time an API call is made about 5 seconds
	// Maybe should make this a bit longer for slower connections?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	todoist := t.NewTodoist(ctx, config)
}
