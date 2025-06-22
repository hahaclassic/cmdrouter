test:
	go test ./...
	
cov:
	mkdir -p ./coverage
	./scripts/cov.sh

cov-v: 
	mkdir -p ./coverage
	./scripts/cov.sh -v

integration:
	go test -v ./integration-tests/

clear:
	rm -rf ./coverage

lint: 
	golangci-lint run