from haproxy:1.8-alpine
run mkdir /var/lib/haproxy
run apk update && \
    apk add openssl
run wget -P /root "https://releases.hashicorp.com/consul-template/0.19.4/consul-template_0.19.4_linux_amd64.tgz"
run cd root && \
    tar xfz consul-template_0.19.4_linux_amd64.tgz && \
    mv consul-template /usr/local/bin
copy haproxy.tpl /etc/haproxy.tpl
copy consul-template.cfg /etc/consul-template.cfg
copy start-or-restart-haproxy.sh /usr/local/bin
run chmod u+x /usr/local/bin/start-or-restart-haproxy.sh
run apk add runit
run mkdir -p /etc/sv/consul-template
copy consul-template-init.sh /etc/sv/consul-template/run
run chmod a+x /etc/sv/consul-template/run
expose 9000
entrypoint ["/sbin/runsvdir","-P","/etc/sv"]
