#!/bin/sh
exec 2>&1
exec /usr/local/bin/consul-template -config /etc/consul-template.cfg
