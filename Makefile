test: ## Run tests 
	cd app && cd interfaces && go test -v

build: ## Run go build
	go mod download
	go mod vendor
	cd app && go build -o scanner

run: ## Start the server
	cd app && ./scanner
