FROM alpine:3.7

RUN apk add --no-cache gcc g++ make cmake gmp-dev
RUN apk add --no-cache doxygen graphviz nodejs
RUN npm install http-server -g

ADD relic /relic

RUN mkdir /relic-build
# RUN cd /relic-build && sh /relic/preset/x64-pbc-128-b12.sh /relic
# https://github.com/relic-toolkit/relic/issues/58
RUN cd /relic-build && cmake -DALLOC=DYNAMIC -DFP_PRIME=381 -DVERBS=off -DRAND=UDEV -DTESTS=1 -DBENCH=0 /relic
RUN make -C /relic-build all doc

EXPOSE 8080
CMD http-server /relic-build/doc/html/
