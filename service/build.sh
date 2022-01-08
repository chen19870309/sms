#!/bin/bash 
cd src
oldGOPATH=$GOPATH
nowPATH=`pwd`
export GOPATH=$oldGOPATH:$nowPATH
echo "GOPATH = $GOPATH"


