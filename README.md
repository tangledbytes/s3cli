# S3CLI
Stupid simple S3 API CLI. All it let's you do is run ***any*** S3 API. You can pass files as well as input which are resolved automatically.

## Why?
Wanted to run `put-bucket-policy` on NooBaa endpoint and `aws s3api` just won't do it. Lookedup other alternatives and none of them satisfied the needs, so here it is.

## How?
```
$ s3cli api --help
Run any S3 API

Usage:
  s3cli api [flags]

Examples:

# Basic usage
s3cli api --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' ListObjects

# File Params example - --file-params are merged with --params after file data expanstion
s3cli api --endpoint http://localhost:9000 --access-key noobaa --secret-key noobaa123 --params '{"Bucket": "test"}' --file-params '{ "Policy": "./policy.json" }'  PutBucketPolicy

# Anonymous usage
s3cli api --endpoint http://localhost:9000 --anon --params '{"Bucket": "test"}' ListObjects


Flags:
      --access-key string    S3 access key
      --anon                 S3 anonymous
      --debug                S3 debug
      --endpoint string      S3 endpoint
      --file-params string   S3 api file params as JSON - gets merged with params after file resolution (default "{}")
  -h, --help                 help for api
  -o, --output string        Output format, one of: [raw color json] (default "json")
      --params string        S3 api params as JSON (default "{}")
      --region string        S3 region (default "us-east-1")
      --secret-key string    S3 secret key
      --skip-ssl             S3 skip ssl
```