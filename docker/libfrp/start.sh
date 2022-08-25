#!/bin/bash

pushd /context
if [[ ! -d libfrp ]]; then
   git clone https://github.com/Darwin-Che/libfrp.git
fi
pushd libfrp
git branch dev
git fetch --all
git reset --hard FETCH_HEAD
git status
make

# the binary file is located at '/context/libfrp/bin'
pushd bin

# the server config file
echo """
[common]
bind_port = $BIND_PORT
vhost_http_port = $HTTP_PORT
subdomain_host = $DOMAIN

dashboard_port = 9500
dashboard_user = admin
dashboard_pwd = admin
""" > server.ini

./frps -c server.ini
