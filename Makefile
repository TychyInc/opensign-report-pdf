.PHONY: samples
samples:
	@echo "Generating sample PDFs..."
	@cd examples/basic && go run main.go
	@echo "Sample PDFs generated in examples/basic/"