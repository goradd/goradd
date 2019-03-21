#!/usr/bin/env bash

# This script is an aid to running the browser test on travis. It is designed to be run from the goradd-test directory.
echo "*** building main"
go build goradd-test
echo "*** starting server"
./goradd-test &
sleep 5
google-chrome-stable --headless --remote-debugging-port=9222 http://localhost:8000/test?all=1 &
#"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless --remote-debugging-port=9222 http://localhost:8000/test?all=1 &
wait %1
r=$?
kill %2
exit $r
