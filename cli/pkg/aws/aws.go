//go:generate go run github.com/utkarsh-pro/s3cli/typereg

package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3" // typereg:s3
)

type AWS struct {
	svc *s3.S3
}

type AWSConfig struct {
	Region    string
	AccessKey string
	SecretKey string
	Anon      bool
	SkipSSL   bool
	Endpoint  string
}

func New(cfg AWSConfig) *AWS {
	config := aws.NewConfig()

	if cfg.Region != "" {
		config = config.WithRegion(cfg.Region)
	}
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""))
	}
	if cfg.Anon {
		config = config.WithCredentials(credentials.AnonymousCredentials)
	}
	if cfg.SkipSSL {
		config = config.WithDisableSSL(true)
	}
	if cfg.Endpoint != "" {
		config = config.WithEndpoint(cfg.Endpoint)
	}

	sess := session.Must(session.NewSession(config))
	return &AWS{
		svc: s3.New(sess),
	}
}

func (a *AWS) RunAny(api string, params map[string]interface{}, fileParams map[string]string) {
	reflect.New(reflect.TypeOf(a.svc).Elem()).Interface()
	method := reflect.ValueOf(a.svc).MethodByName(api)

	i, err := instance(fmt.Sprintf("github.com/aws/aws-sdk-go/service/s3.%sInput", api))
	if err != nil {
		fmt.Println(typeRegistry)
		panic(err)
	}

	merged, err := mergeParams(params, fileParams)
	if err != nil {
		panic(err)
	}

	err = anyToAny(merged, &i)
	if err != nil {
		panic(err)
	}

	invalue := reflect.New(reflect.TypeOf(i)).Elem()
	invalue.Set(reflect.ValueOf(i))

	outputs := method.Call([]reflect.Value{invalue})
	for _, output := range outputs {
		fmt.Println(output)
	}
}

func mergeParams(params map[string]interface{}, fileParams map[string]string) (map[string]interface{}, error) {
	merged := map[string]interface{}{}

	for k, v := range fileParams {
		data, err := os.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		merged[k] = data
	}

	return params, nil
}

func anyToAny(i1, i2 any) error {
	byt, err := json.Marshal(i1)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = json.Unmarshal(byt, i2)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}
