OUTDIR := bin

.PHONY: cli
cli:
	go build -o $(OUTDIR)/s3cli github.com/tangledbytes/s3cli/cli