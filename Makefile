

build:
	go build -o ./backend/bin/moneyKeeper ./backend/cmd/*go

clean:
	rm -rf ./backend/bin/

run:
	go run ./backend/cmd/*.go