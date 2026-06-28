.PHONY: pack dev build build-linux clean tidy gen-dao

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
