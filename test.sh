#!/bin/bash

for pkg in `cat testpackages.txt`
do
    go test $pkg
done
