#!/bin/sh

echo "building..."
go build
if [ $? -eq 0 ]
then
    /home/danny/nwjs-sdk-v0.22.0-linux-x64/nw .
fi
