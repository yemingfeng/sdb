FROM ubuntu:18.04

COPY sdb sdb

ENTRYPOINT "./sdb"