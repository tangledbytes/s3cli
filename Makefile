OUTDIR := bin

.PHONY: cli
cli:
	go build -o $(OUTDIR)/cli github.com/utkarsh-pro/s3cli/cli