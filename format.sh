#!/bin/sh
for f in $(find . -type f -name "*.go"); do 
	gofmt -w $f; 
done
