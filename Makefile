<<<<<<< HEAD
hello: 
	echo "Hello"
run:
	go run cmd/main.go
=======
bins: esm esmd

esm:
	go build ./cmd/esm

esmd:
	go build ./cmd/esmd

tester:
	go test ./test -run ''

clean:
	rm -f esm esmd
>>>>>>> df02bfd76a1796128d65681a9d487de5e9d0fd22
