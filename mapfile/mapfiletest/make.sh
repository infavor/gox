#!/bin/sh

go build -o appendonly appendonly.go
go build -o appendcheck appendcheck.go
go build -o depthcheck depthcheck.go