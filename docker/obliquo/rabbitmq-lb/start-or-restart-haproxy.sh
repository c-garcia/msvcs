#!/bin/sh

HAPROXY=/usr/local/sbin/haproxy
HAPROXY_CFG=/etc/haproxy.cfg
HAPROXY_PID=/var/run/haproxy.pid
HAPROXY_SOCK=/var/run/haproxy-admin.sock

if [ -S ${HAPROXY_SOCK} ]; then
    ${HAPROXY} -f ${HAPROXY_CFG} \
               -p ${HAPROXY_PID} \
               -x ${HAPROXY_SOCK} \
               -sf $(pidof haproxy)
else 
    ${HAPROXY} -f ${HAPROXY_CFG}
fi
