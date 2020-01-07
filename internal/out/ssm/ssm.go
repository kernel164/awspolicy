package ssm

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSM struct {
	ssm *ssm.SSM
}

func newSSM() *SSM {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	return &SSM{ssm: ssm.New(sess)}
}

func (s *SSM) Get(key string) (string, error) {
	key, encrypted := s.normalize(key)
	fmt.Printf("==> SSM: Retrieving Parameter %s\n", key)
	param, err := s.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: &encrypted,
	})
	if err != nil {
		return "", err
	}
	return *param.Parameter.Value, nil
}

func (s *SSM) Set(key string, value string) error {
	key, _ = s.normalize(key)
	fmt.Printf("==> SSM: Updating Parameter %s\n", key)
	_, err := s.ssm.PutParameter(&ssm.PutParameterInput{
		Name:      &key,
		Value:     &value,
		Overwrite: aws.Bool(true),
		Type:      aws.String("String"),
	})
	if err == nil {
		fmt.Printf("==> SSM: Updated Parameter %s\n", key)
	}
	return err
}

func (s *SSM) normalize(key string) (string, bool) {
	if strings.HasPrefix(key, "ssm:") {
		key = key[4:]
	}
	return key, false
}
