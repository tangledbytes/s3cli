package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/s3cli/cli/pkg/aws"
	"github.com/utkarsh-pro/s3cli/cli/pkg/printer"
	"github.com/utkarsh-pro/s3cli/cli/pkg/restrictedflag"
	"github.com/utkarsh-pro/s3cli/cli/pkg/utils"
)

const apiExample = `
# Basic usage
s3cli api --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' ListObjects

# File Params example - --file-params are merged with --params after file data expanstion
s3cli api --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' --file-params '{ "Policy": "./policy.json" }'  PutBucketPolicy

# Anonymous usage
s3cli api --endpoint http://localhost:9000 --anon --params '{"Bucket": "test"}' ListObjects
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
	debug      bool
	outputType = restrictedflag.New("json", "raw", "color", "json")
)

var (
	runner           *aws.AWS
	parsedParams     map[string]interface{}
	parsedFileParams map[string]string
)

func generateValidArgs() []string {
	disallowed := []string{
		"Parquet",
		"JSON",
		"CSV",
	}

	valid := []string{}
	types := aws.GetTypeRegistry()
	for _, typ := range types {
		name := typ.Name()
		if strings.HasSuffix(name, "Input") &&
			!utils.ContainsAny(
				disallowed,
				[]string{name},
				func(v1, v2 string) bool { return v1 == strings.TrimSuffix(v2, "Input") },
			) {
			valid = append(valid, strings.TrimSuffix(name, "Input"))
		}
	}

	return valid
}

// RootCmd represents the base command when called without any subcommands
var ApiCmd = &cobra.Command{
	Use:       "api",
	Short:     "Run any S3 API",
	Example:   apiExample,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: generateValidArgs(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		runner = aws.New(aws.AWSConfig{
			Region:    region,
			Endpoint:  endpoint,
			AccessKey: access,
			SecretKey: secret,
			Anon:      anon,
			SkipSSL:   skipSSL,
			Debug:     debug,
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

func init() {
	ApiCmd.Flags().StringVar(&endpoint, "endpoint", "", "S3 endpoint")
	ApiCmd.Flags().StringVar(&region, "region", "us-east-1", "S3 region")
	ApiCmd.Flags().StringVar(&access, "access-key", "", "S3 access key")
	ApiCmd.Flags().StringVar(&secret, "secret-key", "", "S3 secret key")
	ApiCmd.Flags().BoolVar(&anon, "anon", false, "S3 anonymous")
	ApiCmd.Flags().BoolVar(&skipSSL, "skip-ssl", false, "S3 skip ssl")
	ApiCmd.Flags().StringVar(&params, "params", "{}", "S3 api params as JSON")
	ApiCmd.Flags().StringVar(&fileParams, "file-params", "{}", "S3 api file params as JSON - gets merged with params after file resolution")
	ApiCmd.Flags().BoolVar(&debug, "debug", false, "S3 debug")

	ApiCmd.Flags().VarP(outputType, "output", "o", fmt.Sprintf("Output format, one of: %s", outputType.Allowed()))

	ApiCmd.MarkFlagRequired("endpoint")
	ApiCmd.MarkFlagsRequiredTogether("access-key", "secret-key")
	ApiCmd.MarkFlagsMutuallyExclusive("anon", "access-key")
	ApiCmd.MarkFlagsMutuallyExclusive("anon", "secret-key")
}
