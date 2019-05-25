#!/usr/bin/env bash

set -uo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

function format_dir() {
	echo "Formating $1"
	pushd $1 > /dev/null
	go fmt
	popd > /dev/null
}

format_dir $DIR
format_dir $DIR/overlay
format_dir $DIR/underlay
format_dir $DIR/capacity
format_dir $DIR/structs
format_dir $DIR/sdk/go
format_dir $DIR/logs
format_dir $DIR/interfaces
format_dir $DIR/examples
format_dir $DIR/events
format_dir $DIR/config
