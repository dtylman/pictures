#!/bin/sh

echo "building..."
go build
if [ $? -eq 0 ]
then
    /home/danny/src/nwjs/nwjs-sdk-v0.21.3-linux-x64/nw .
fi
