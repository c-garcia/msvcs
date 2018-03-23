global
  daemon
  chroot /var/lib/haproxy
  pidfile /var/run/haproxy.pid
  maxconn 4096
  stats socket /var/run/haproxy-admin.sock mode 660 level admin expose-fd listeners

defaults
  mode tcp
  maxconn 10000
  timeout connect 5s
  timeout client 100s
  timeout server 100s

frontend rabbitmq
  bind *:5672
  mode tcp
  timeout client 3h
  option clitcpka
  default_backend be_rabbitmq

frontend stats
  bind *:9000
  mode http
  stats enable
  stats refresh 30s
  stats uri /
  stats admin if TRUE

backend be_rabbitmq
  balance leastconn
  {{range service "rabbitmq"}}
  server {{.ID}} {{.Address}}:{{.Port}} check
  {{end}}