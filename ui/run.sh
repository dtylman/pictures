#!/bin/sh

cd ../cmd/pc
echo "buidling..."
go build
if [ $? -eq 0 ]
then
    cd ../../ui
    /home/danny/src/nwjs/nwjs-sdk-v0.21.3-linux-x64/nw .
fi
cd ../../ui