#!/bin/bash

. config.inc.sh

bin=./flightradar-telegram-bot
if [ ! -x "$bin" ]; then
	bin="go run *.go"
fi

BOT_TOKEN=$BOT_TOKEN \
ALLOWED_USERIDS=$ALLOWED_USERIDS \
ADMIN_USERIDS=$ADMIN_USERIDS \
ALLOWED_GROUPIDS=$ALLOWED_GROUPIDS \
$bin $*
