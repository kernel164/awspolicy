package out

import (
	"encoding/json"
	"fmt"
)

type stdout struct {
}

func New() Output {
	return &stdout{}
}

func (o *stdout) Get() (*PolicyDocument, error) {
	return &PolicyDocument{Version: "2012-10-17", Statement: []*PolicyStatement{&PolicyStatement{Action: []string{}, Effect: "Allow", Resource: "*"}}}, nil
}

func (o *stdout) Set(doc *PolicyDocument) error {
	fmt.Println("==> Policy Document:")
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(string(data))
	fmt.Println()
	return nil
}
