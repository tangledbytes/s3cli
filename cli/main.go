package main

import (
	"flag"
	"fmt"

	"github.com/utkarsh-pro/s3cli/cli/pkg/aws"
)

type A struct{}

func setupFlag() map[string]interface{} {
	flags := make(map[string]interface{})

	flags["endpoint"] = flag.String("endpoint", "", "S3 endpoint")
	flags["region"] = flag.String("region", "", "S3 region")
	flags["access-key"] = flag.String("access-key", "", "S3 access key")
	flags["secret-key"] = flag.String("secret-key", "", "S3 secret key")
	flags["anon"] = flag.Bool("anon", false, "S3 anonymous")
	flags["skip-ssl"] = flag.Bool("skip-ssl", false, "S3 skip ssl")
	flags["api"] = flag.String("api", "", "S3 api name")
	flags["params"] = flag.String("params", "", "S3 api params")
	flags["file-params"] = flag.String("file-params", "", "S3 api file params - gets merged with params after file resolution")

	flag.Parse()

	return flags
}

func ensureRequired(flags map[string]interface{}, required []string) error {
	for _, r := range required {
		if _, ok := flags[r]; !ok {
			return fmt.Errorf("%s flag is required", r)
		}
	}

	return nil
}

func main() {
	flags := setupFlag()

	if err := ensureRequired(flags, []string{"api", "endpoint"}); err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		return
	}

	aws.
		New(aws.AWSConfig{
			Region: flags["region"].(string),
		}).
		RunAny("ListBuckets", map[string]interface{}{}, map[string]string{})
}
