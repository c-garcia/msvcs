from rabbitmq:3-management-alpine
run rabbitmq-plugins --offline enable rabbitmq_peer_discovery_consul
copy consul.conf /tmp
run cat /tmp/consul.conf >> /etc/rabbitmq/rabbitmq.conf
run rm /tmp/consul.conf

