all: goweblua

lualib:
	cd c && make

goweblua: lualib
	CGO_ENABLED=1 go build -o goweblua cmd/main.go