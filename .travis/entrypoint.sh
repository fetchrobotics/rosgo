#!/bin/bash
source /opt/ros/melodic/setup.bash
export PATH=$PWD/bin:/usr/local/go/bin:$PATH
export GOPATH=$PWD:/usr/local/go

roscore &
go install github.com/fetchrobotics/rosgo/gengo
go generate github.com/fetchrobotics/rosgo/tests
go test github.com/fetchrobotics/rosgo/xmlrpc
go test github.com/fetchrobotics/rosgo/ros
go test github.com/fetchrobotics/rosgo/tests/...
