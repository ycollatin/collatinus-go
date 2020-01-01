defaut:
	go build -o bin/collatinus
darwin:
	env GOOS=darwin GOARCH=amd64 go build -o mac/collatinus
w:
	env GOOS=windows GOARCH=amd64 go build -o win/collatinus.exe
edit:
	vim -c "b: coll.go" *.go
