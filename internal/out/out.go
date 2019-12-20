package out

type PolicyDocument struct {
	Version   string
	Statement []*PolicyStatement
}

type PolicyStatement struct {
	Action   []string
	Effect   string
	Resource string
}

type Output interface {
	Get() (*PolicyDocument, error)
	Set(*PolicyDocument) error
}

type Engine struct {
}
