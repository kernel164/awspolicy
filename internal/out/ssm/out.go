package ssm

import (
	"encoding/json"
	"strings"

	"github.com/kernel164/awspolicy/internal/out"
)

type ssmOutput struct {
	path string
	ssm  *SSM
}

func New(path string) out.Output {
	return &ssmOutput{path: path, ssm: newSSM()}
}

func (o *ssmOutput) Get() (*out.PolicyDocument, error) {
	val, err := o.ssm.Get(o.path)
	if err != nil {
		if !strings.Contains(err.Error(), "ParameterNotFound") {
			return nil, err
		}
	}
	if val == "" {
		return &out.PolicyDocument{Version: "2012-10-17", Statement: []*out.PolicyStatement{&out.PolicyStatement{Action: []string{}, Effect: "Allow", Resource: "*"}}}, nil
	}
	doc := &out.PolicyDocument{}
	err = json.Unmarshal([]byte(val), doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (o *ssmOutput) Set(doc *out.PolicyDocument) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return o.ssm.Set(o.path, string(data))
}
