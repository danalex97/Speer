#!/bin/bash
if [ -z "$GOPATH" ]
then
      echo "\$GOPATH was not set. Please set it before running the installer."
      exit 1
fi

mkdir -p $GOPATH/src/github.com/danalex97/Speer

pushd $GOPATH/src/github.com/danalex97 > /dev/null

if [ -z "$(ls -A Speer)" ]; then
    rm -rf Speer
    git clone https://github.com/danalex97/Speer.git
else
    pushd Speer > /dev/null
    git pull
    popd > /dev/null
fi

pushd Speer > /dev/null

chmod +x speer.sh
sudo cp speer.sh /usr/bin/speer
echo "Speer installed successfully."

popd > /dev/null

popd > /dev/null
