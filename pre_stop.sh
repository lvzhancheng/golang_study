#!/bin/bash
echo "stop http_server"
ps -ef | grep http_server | xargs kill -1
