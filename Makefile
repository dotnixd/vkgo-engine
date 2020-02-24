SRC=$(filter-out plugins/basic/main.go, $(wildcard plugins/*/*.go))
OBJ=$(SRC:.go=.so)

vkgo: $(OBJ)
	go build

%.so: %.go
	go build -buildmode=plugin -o $@ $<

dependencies:
	go get github.com/valyala/fastjson
	go get github.com/jmoiron/sqlx
	go get github.com/go-sql-driver/mysql

clean:
	rm $(OBJ) -f
	rm $(notdir $(CURDIR)) -f
	rm -rf release

release: clean vkgo
	mkdir release/ -p
	cp $(notdir $(CURDIR)) release/
	cp errors.json release/
	cp plugins release/ -R
	rm -rf release/plugins/*/*.go release/plugins/basic