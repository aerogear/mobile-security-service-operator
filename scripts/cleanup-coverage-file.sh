#!/bin/sh

# Do a check to see if the user is on MacOS or Linux
if [[ "$OSTYPE" == "darwin"* ]] || [[ "$OSTYPE" == "freebsd"* ]]; then
	## Remove files ending in _mock.go
	sed -i '' '/_mock.go/d' ./coverage-all.out
else
	## Remove files ending in _mock.go
	sed -i '/_mock.go/d' ./coverage-all.out
fi