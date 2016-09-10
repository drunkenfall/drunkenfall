#!/usr/bin/env bash

cd $(dirname $(readlink -f $0))

# Check that we have the things that we need to run this
error=false
for x in go npm nginx; do
    if ! type $x &> /dev/null ; then
        echo "$x not found. Please install."
        error=true
    fi
done

if $error ; then
    exit 1
fi

function start_nginx() {
    mkdir -p logs
    # TODO: Make sure we can run this without sudo
    sudo nginx -p . -c conf/nginx.conf
}

function start_api() {
    # This file contains the environment variables that configures the
    # Facebook application, and since parts of those are secret, they are not
    # found in this repository.

    f="df_fb.env"
    if [[ -f $f ]]; then
        source $f
    fi

    go build -v || exit $?
    ./drunkenfall
}

function start_npm() {
    cd js/
    if [[ ! -d "node_modules" ]]; then
        echo "No modules found - will download."
        npm install
    fi

    npm run dev
}

case "$1" in
    api)
        start_api
        ;;
    nginx)
        start_nginx
        ;;
    npm)
        start_npm
        ;;
    *)
        echo "usage: dev.sh [towerfall|nginx|npm]"
        exit 1
esac
