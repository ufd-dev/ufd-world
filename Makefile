dev:
	go run . -i
deploy:
	go build .
	scp ufd-world admin@ufd.world:~/