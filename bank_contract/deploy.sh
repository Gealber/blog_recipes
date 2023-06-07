#! /bin/bash

date
cd scripts/ && sh compile.sh && cd ../ && go run main.go
