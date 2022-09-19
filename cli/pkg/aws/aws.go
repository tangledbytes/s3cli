//go:generate go run github.com/utkarsh-pro/s3cli/typereg

package aws

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3" // typereg:s3
	"github.com/utkarsh-pro/s3cli/cli/pkg/utils"
)

// AWS is a wrapper around the AWS SDK having helper
// functions to run any API.
type AWS struct {
	svc   *s3.S3
	debug bool
}

// AWSConfig is the configuration for AWS.
type AWSConfig struct {
	Region    string
	AccessKey string
	SecretKey string
	Anon      bool
	SkipSSL   bool
	Endpoint  string
	Debug     bool
}

// New consumes a config and returns a new AWS instance.
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
		config.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		})
		config = config.WithDisableSSL(true)
	}
	if cfg.Endpoint != "" {
		config = config.WithEndpoint(cfg.Endpoint)
	}

	sess := session.Must(session.NewSession(config))
	return &AWS{
		svc:   s3.New(sess),
		debug: cfg.Debug,
	}
}

// RunAny takes an API name and a map of params and runs the API.
func (a *AWS) RunAny(api string, params map[string]interface{}, fileParams map[string]string) ([]interface{}, error) {
	method := reflect.ValueOf(a.svc).MethodByName(api)

	i, err := NewInstance(fmt.Sprintf("github.com/aws/aws-sdk-go/service/s3.%sInput", api))
	if err != nil {
		return nil, fmt.Errorf("failed to get input instance: %w", err)
	}

	merged, err := mergeParams(params, fileParams)
	if err != nil {
		return nil, fmt.Errorf("failed to merge params: %w", err)
	}

	err = utils.AnyToAny(merged, i)
	if err != nil {
		return nil, fmt.Errorf("failed to convert params to input: %w", err)
	}

	if a.debug {
		return []interface{}{i}, nil
	}

	invalue := reflect.New(reflect.TypeOf(i)).Elem()
	invalue.Set(reflect.ValueOf(i))

	outputs, err := utils.ValueSliceToInterfaceSlice(method.Call([]reflect.Value{invalue}))
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

// ParseParams takes params as string and returns a map of params.
func (a *AWS) ParseParams(params string) (map[string]interface{}, error) {
	return utils.ParseJSONToMapStringInterface(params)
}

// ParseFileParams takes file params as string and returns a map of file params.
func (a *AWS) ParseFileParams(params string) (map[string]string, error) {
	return utils.ParseJSONToMapStringString(params)
}

// mergeParams takes a params and fileParams and fills fileParams with the file contents and
// returns a merged map.
func mergeParams(params map[string]interface{}, fileParams map[string]string) (map[string]interface{}, error) {
	merged := map[string]interface{}{}

	for k, v := range fileParams {
		data, err := os.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		merged[k] = string(data)
	}

	for k, v := range params {
		merged[k] = v
	}

	return merged, nil
}
