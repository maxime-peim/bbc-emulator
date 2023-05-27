SOURCE=bbc

test:
	go test $(SOURCE)/...

run:
	go run $(SOURCE)/main.go