#!/bin/bash
if [ -z "$GOPATH" ]
then
      echo "\$GOPATH was not set. Please set it before running the installer."
	  exit 1
fi

SPEER=github.com/danalex97/Speer

# Download Speer
go get -d $SPEER

pushd $GOPATH/src/$SPEER >> /usr/null

chmod +x speer.sh
cp speer.sh /usr/bin/speer
echo "Speer installed successfully."

popd >> /usr/null
