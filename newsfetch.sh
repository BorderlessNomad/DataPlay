#!/bin/bash

echo 'BUILDING GOGRAM'
oldgo=$GOPATH
if [[ "$OSTYPE" == "msys" ]]; then
	GOPATH=$oldgo";"$(pwd -W)
else
	GOPATH=$oldgo:$(pwd)
fi
export GOPATH
project=newsfetch
go get -v $project
go install -v $project
cd bin
./$project $@
export GOPATH=$oldgo
