#!/usr/bin/env bash

TEST_SOCK=./foobar.sock

trap 'unlink ${TEST_SOCK}' INT

while true; do date; sleep 1; done | nc -lkU ${TEST_SOCK} -

