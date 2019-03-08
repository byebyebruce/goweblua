all: goweblua

goweblua:
	cd c && make
	go build