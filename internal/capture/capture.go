package capture

type Params struct {
	Args []string
}

type Results struct {
	Actions map[string]struct{}
}

type Capture interface {
	Run(*Params) (*Results, error)
}
