package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kernel164/awspolicy/internal/out"
)

type fileOutput struct {
	path string
}

func New(path string) out.Output {
	return &fileOutput{path: path}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (o *fileOutput) Get() (*out.PolicyDocument, error) {
	policy := &out.PolicyDocument{Version: "2012-10-17", Statement: []*out.PolicyStatement{&out.PolicyStatement{Action: []string{}, Effect: "Allow", Resource: "*"}}}
	if fileExists(o.path) {
		content, err := ioutil.ReadFile(o.path)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, policy)
		if err != nil {
			return nil, err
		}
	} else {
		err := ioutil.WriteFile(o.path, []byte(""), 0644)
		if err != nil {
			return nil, err
		}
	}
	return policy, nil
}

func (o *fileOutput) Set(doc *out.PolicyDocument) error {
	fmt.Println("==> OutFile:", o.path)
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(o.path, data, 0644)
}
