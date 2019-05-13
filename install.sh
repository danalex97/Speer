#!/bin/bash
if [ -z "$GOPATH" ]
then
      echo "\$GOPATH was not set. Please set it before running the installer."
      exit 1
fi

mkdir -p $GOPATH/src/danalex97

pushd $GOPATH/src >> /usr/null
if [ -z "$(ls -A Speer)" ]; then
    rm -rf Speer
    git clone https://github.com/danalex97/Speer.git
else
    pushd Speer >> /usr/null
    git pull
    popd >> /usr/null
fi

pushd Speer >> /usr/null

chmod +x speer.sh
sudo cp speer.sh /usr/bin/speer
echo "Speer installed successfully."

popd >> /usr/null

popd >> /usr/null
