#!/bin/bash
if [ -z "$GOPATH" ]
then
      echo "\$GOPATH was not set. Please set it before running Speer."
	  exit 1
fi

go run $GOPATH/src/github.com/danalex97/Speer/speer.go $@
