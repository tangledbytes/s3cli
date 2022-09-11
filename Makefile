OUTDIR := bin

.PHONY: cli
cli:
	go build -o $(OUTDIR)/s3cli github.com/utkarsh-pro/s3cli/cli