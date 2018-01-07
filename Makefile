.PHONY: relic
relic:
	rm -rf bls12/build && mkdir bls12/build
	cd bls12/build && cmake -DALLOC=DYNAMIC -DFP_PRIME=381 \
		-DSHLIB=off -DSTLIB=on -DRAND=UDEV -DTESTS=1 -DBENCH=0 \
		-DCOMP="-O3 -funroll-loops -Wno-unused-function" ../relic
	make -C bls12/build
	make -C bls12/build test

.PHONY: docs
docs:
	docker build -t relic bls12
	docker run -it -p 8080:8080 relic
