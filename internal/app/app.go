package app

// App is a struct containing all application services
type App struct {
}

// Takes all services as params and returns a pointer to App
func New() *App {
	return &App{}
}
