#!/bin/bash

# Golang dependencies
echo "Getting Golang dependencies..."
go get github.com/gorilla/mux

# Javascript dependencies
echo "Downloading JS libraries..."
pushd static/libs > /dev/null
wget -nc https://unpkg.com/react@16/umd/react.development.js &> /dev/null
wget -nc https://unpkg.com/react-dom@16/umd/react-dom.development.js &> /dev/null
wget -nc https://unpkg.com/babel-standalone@6/babel.min.js &> /dev/null
wget -nc https://d3js.org/d3.v5.min.js &> /dev/null
popd > /dev/null
