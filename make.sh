#!/bin/bash
protoName=$2
case $1 in
    'go')
        libraryDir="/Users/Mr.Zhou/Project/golang/Libraries/src/local/"
        rm -rf ${libraryDir}${protoName}
        mkdir ${libraryDir}${protoName}
        protoc -I ./go-server/ --proto_path=./ --go_out=plugins=grpc:${libraryDir}${protoName} ./meta/$protoName.proto
    ;;
    'php')
        phpPluginDir="/Users/Mr.Zhou/Project/PHP/grpc-master/bins/opt/grpc_php_plugin"
        phpDir="./php-client/src/proto/"
        protoc --proto_path=./ --php_out=$phpDir --grpc_out=$phpDir --plugin=protoc-gen-grpc=$phpPluginDir ./meta/$protoName.proto
    ;;
esac

