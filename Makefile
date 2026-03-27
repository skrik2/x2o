.PHONY: build clean



build: build-backend
	@echo "Full build completed!"

# build-backend
VERSION ?= test
LDFLAGS := -s -w -X github.com/skrik2/x2o.Version=$(VERSION)
build-backend:
	@echo "Building x2o backend..."
	go build -ldflags "$(LDFLAGS)" -o x2o ./x2ocmd/cmd
	@echo "Backend build completed!"

clean:
	- rm -rf x2o