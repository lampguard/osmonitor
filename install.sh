#!/bin/sh

sudo curl https://github.com/lampguard/osmonitor/archive/refs/tags/latest.tar.gz -o osmonitor.tar.gz

sudo tar -xvf osmonitor.tar.gz
sudo cp osmonitor/RELEASE/osmonitor /usr/local/osmonitor

sudo useradd -r -s /bin/bash lampguard

sudo cp osmonitor/RELEASE/osmonitor.service /lib/systemd/system/osmonitor.service