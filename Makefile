.PHONY: pack dev build build-linux clean tidy gen-dao server-build server-dev docker-build docker-push

# pack embeds manifest resources into the binary with gf pack
pack:
	gf pack manifest/migrations packed/packed.go -y -n packed

# generate DAO/model code from current database schema
gen-dao:
	gf gen dao

dev: pack
	wails3 dev

build: pack
	wails3 build

build-linux: pack
	wails3 build -tags "webkit2_41"

clean:
	rm -rf build/ bin/ frontend/dist/ frontend/node_modules/ app.log data/
	rm -f packed/packed.go

tidy:
	go mod tidy
	cd frontend && npm install

server-build:
	cd server && go mod tidy && go build -o ../bin/ydsterm-server .

server-dev:
	cd server && go run .

server-gen-dao:
	cd server && gf gen dao

docker-build:
	@if [ -z "$(REGISTRY)" ] || [ -z "$(VERSION)" ]; then \
		echo "Usage: make docker-build REGISTRY=<registry/namespace> VERSION=<version>"; \
		echo "Example: make docker-build REGISTRY=gitea.xxx.com/namespace VERSION=v1.0.0"; \
		exit 1; \
	fi
	cd server && docker build -t $(REGISTRY)/ydsterm-server:$(VERSION) .
	@echo "Image built: $(REGISTRY)/ydsterm-server:$(VERSION)"

docker-push:
	@if [ -z "$(REGISTRY)" ] || [ -z "$(VERSION)" ]; then \
		echo "Usage: make docker-push REGISTRY=<registry/namespace> VERSION=<version>"; \
		exit 1; \
	fi
	docker push $(REGISTRY)/ydsterm-server:$(VERSION)
