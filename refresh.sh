#!/bin/sh

rm ./routers/commentsRouter_github_com_lfq7413_tomato_controllers.go
echo "remove commentsRouter success"
mv main.go.t main.go
go run main.go
echo "run main.go success"
mv main.go main.go.t
echo "refresh success"
rm lastupdate.tmp
