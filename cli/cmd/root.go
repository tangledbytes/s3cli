package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/s3cli/cli/pkg/aws"
	"github.com/utkarsh-pro/s3cli/cli/pkg/printer"
	"github.com/utkarsh-pro/s3cli/cli/pkg/restrictedflag"
)

const example = `
# Basic usage
s3cli --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' ListObjects

# File Params example - --file-params are merged with --params after file data expanstion
s3cli --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' --file-params '{ "Policy": "./policy.json" }'  PutBucketPolicy

# Anonymous usage
s3cli --endpoint http://localhost:9000 --anon --params '{"Bucket": "test"}' ListObjects
`

var (
	endpoint   string
	region     string
	access     string
	secret     string
	anon       bool
	skipSSL    bool
	params     string
	fileParams string
	outputType = restrictedflag.New("json", "raw", "color", "json")
)

var (
	runner           *aws.AWS
	parsedParams     map[string]interface{}
	parsedFileParams map[string]string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "s3cli",
	Short:   "S3CLI is a stupid simple CLI for S3",
	Example: example,
	Args:    cobra.ExactArgs(1),
	Version: "0.0.1",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		runner = aws.New(aws.AWSConfig{
			Region:    region,
			Endpoint:  endpoint,
			AccessKey: access,
			SecretKey: secret,
			Anon:      anon,
			SkipSSL:   skipSSL,
		})

		parsedParams, err = runner.ParseParams(params)
		if err != nil {
			return fmt.Errorf("failed to parse params: %w", err)
		}

		parsedFileParams, err = runner.ParseFileParams(fileParams)
		if err != nil {
			return fmt.Errorf("failed to parse file params: %w", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		output, err := runner.RunAny(args[0], parsedParams, parsedFileParams)
		if err != nil {
			fmt.Println(err)
			return
		}

		printer.Print(output, outputType.Get() == "color", outputType.Get() == "raw")
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().StringVar(&endpoint, "endpoint", "", "S3 endpoint")
	RootCmd.Flags().StringVar(&region, "region", "us-east-1", "S3 region")
	RootCmd.Flags().StringVar(&access, "access-key", "", "S3 access key")
	RootCmd.Flags().StringVar(&secret, "secret-key", "", "S3 secret key")
	RootCmd.Flags().BoolVar(&anon, "anon", false, "S3 anonymous")
	RootCmd.Flags().BoolVar(&skipSSL, "skip-ssl", false, "S3 skip ssl")
	RootCmd.Flags().StringVar(&params, "params", "{}", "S3 api params as JSON")
	RootCmd.Flags().StringVar(&fileParams, "file-params", "{}", "S3 api file params as JSON - gets merged with params after file resolution")

	RootCmd.Flags().VarP(outputType, "output", "o", fmt.Sprintf("Output format, one of: %s", outputType.Allowed()))

	RootCmd.MarkFlagRequired("endpoint")
}
