bins: esm esmd

esm:
	go build ./cmd/esm

esmd:
	go build ./cmd/esmd

tester:
	go test ./test -run ''

clean:
	rm -f esm esmd
