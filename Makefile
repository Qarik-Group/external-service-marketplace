bins: esm esmd

esm:
	go build ./cmd/esm

esmd:
	go build ./cmd/esmd

clean:
	rm -f esm esmd
