#!/usr/bin/env bash

cd ../../../deploy
# Put your sftp command below, filling in the values. You should configure ssh to connect to your server without using
# a password. You also will need to make this file editable. If on Windows, do it another way
sftp {user}@{server.com}:{/server/path/to/deployment/bin} <<EOF
put {application-in-deploy}
EOF