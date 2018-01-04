.PHONY: relic
relic:
	rm -rf relic/build && mkdir relic/build
	cd relic/build && cmake -DFP_PRIME=381 -DVERBS=off -DBENCH=off ../relic
	make -C relic/build
	make -C relic/build test

.PHONY: docs
docs:
	docker build -t relic relic
	docker run -it -p 8080:8080 relic
