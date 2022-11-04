build:
	@echo 'Build start'
	go build -v -o ./bin/stats ./cmd/stats
	@echo 'Build end'