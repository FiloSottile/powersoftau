.PHONY: relic
relic:
	rm -rf relic/build && mkdir relic/build
	# https://github.com/relic-toolkit/relic/issues/61
	cd relic/build && cmake -DALLOC=DYNAMIC -DFP_PRIME=381 \
		-DSHLIB=off -DSTLIB=on -DRAND=UDEV -DTESTS=1 -DBENCH=0 \
		-DCOMP="-O3 -funroll-loops -Wno-unused-function" \
		-DFP_METHD="BASIC;COMBA;COMBA;BASIC;BASIC;SLIDE" \
		../relic
	make -C relic/build
	make -C relic/build test

.PHONY: docs
docs:
	docker build -t relic relic
	docker run -it -p 8080:8080 relic
