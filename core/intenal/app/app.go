package app

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	return nil
}
