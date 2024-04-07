#!/bin/sh

curl https://github.com/lampguard/osmonitor/archive/refs/tags/latest.tar.gz -L -o osmonitor.tar.gz

tar -xvf osmonitor.tar.gz
cp osmonitor-latest/RELEASE/osmonitor /usr/local/osmonitor

useradd -r -s /bin/bash lampguard

cp osmonitor-latest/RELEASE/osmonitor.service /lib/systemd/system/osmonitor.service