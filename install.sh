#!/bin/sh

rm osmonitor_$(uname -m) osmonitor.service

wget https://github.com/lampguard/osmonitor/raw/main/RELEASE/osmonitor_$(uname -m)
wget https://github.com/lampguard/osmonitor/raw/main/RELEASE/osmonitor.service

rm /lib/systemd/system/osmonitor.service

chmod +x osmonitor_$(uname -m)

cp osmonitor_$(uname -m) /usr/local/bin/osmonitor && echo "Binary installed"
cp osmonitor.service /lib/systemd/system/ && echo "Service configured"

echo "Run 'osmonitor --login' to setup your account"
echo "Then run service osmonitor start to start the service"
echo "Happy monitoring!!"