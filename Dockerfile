FROM ubuntu:18.04

COPY sdb sdb

COPY configs/config.yml configs/config.yml

ENTRYPOINT "./sdb" $0 $@