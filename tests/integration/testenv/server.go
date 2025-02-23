package testenv

import (
	"github.com/ayo-awe/go-backend-starter/internal/app"
	"github.com/ayo-awe/go-backend-starter/internal/server"
)

func (t *TestEnvironment) setupServer() error {
	t.Server = server.New(app.New(), t.Logger)

	return nil
}
