all: goweblua

goweblua:
	cd c && make
	go build -o goweblua cmd/main.go