default: build

clean:
	rm -rf public frontend/dist git-http-server

build_ui: clean
	cd frontend && grunt build && cp -r dist/ ../public

build: build_ui
	go build git-http-server.go


