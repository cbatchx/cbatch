#!/bin/bash
yum -y install epel-release

echo $HOSTNAME > /etc/torque/server_name

yum -y install torque-scheduler torque-server torque-client

# Fixing stuff...
dd if=/dev/urandom of=/etc/munge/munge.key bs=1 count=1024
chmod 400 /etc/munge/munge.key
chown munge:munge /etc/munge/munge.key
systemctl enable munge.service
systemctl start munge.service
systemctl enable trqauthd.service
systemctl start trqauthd.service

yes | pbs_server -t create
qmgr -c "set server acl_hosts=master"
qmgr -c "set server scheduling=true"
qmgr -c "create queue batch queue_type=execution"
qmgr -c "set queue batch started=true"
qmgr -c "set queue batch enabled=true"
qmgr -c "set queue batch resources_default.nodes=1"
qmgr -c "set queue batch resources_default.walltime=3600"
qmgr -c "set server default_queue=batch"
qmgr -c "set server keep_completed = 10"
qterm

# Build up the nodes file
: > /var/lib/torque/server_priv/nodes
for i in `seq 1 $1`; do echo  "slave$i np=1" >> /var/lib/torque/server_priv/nodes; done

# Enable and start services
systemctl enable pbs_server.service
systemctl start pbs_server.service
systemctl enable pbs_sched.service
systemctl start pbs_sched.service

# Configure NFS
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server

# This is only a test setup, this is not secure
echo "/home *(rw,sync,no_root_squash,no_all_squash)" > /etc/exports
systemctl restart nfs-server
