# S3CLI
Stupid simple S3 API CLI. All it let's you do is run ***any*** S3 API. You can pass files as well as input which are resolved automatically.

## Why?
Wanted to run `put-bucket-policy` on NooBaa endpoint and `aws s3api` just won't do it. Lookedup other alternatives and none of them satisfied the needs, so here it is.

## How?
```
$ s3cli --help
Usage: s3cli [options]

Stupid simple Amazon S3 API CLI

Options:
  -V, --version            output the version number
  --endpoint <endpoint>    endpoint
  --api <apiName>          api name
  --accessKey <accessKey>  access key
  --secretKey <secretKey>  secret key
  --params <params>        params
  --fp <fp>                fp is file parameter
  --tls <tls>              tls (default: false)
  --anon                   anonymous (default: false)
  --skip-ssl               skip ssl (default: false)
  --debug                  debug
  -h, --help               display help for command
```