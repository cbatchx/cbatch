#!/bin/bash
yum -y install wget

# Install Torque server
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1.2/torque-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1.2/torque-server-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-server.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1.2/torque-scheduler-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-scheduler.rpm

yum -y --nogpgcheck localinstall torque.rpm torque-server.rpm torque-scheduler.rpm

echo $HOSTNAME > /var/spool/torque/server_name
mkdir -p /var/spool/torque/checkpoint/ # Folder required by pbs_server

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
: > /var/spool/torque/server_priv/nodes
for i in `seq 1 $1`; do echo  "slave$i np=1" >> /var/spool/torque/server_priv/nodes; done

# Enable and start services
systemctl enable trqauthd.service
systemctl start trqauthd.service
systemctl enable pbs_server.service
systemctl start pbs_server.service
systemctl enable pbs_sched.service
systemctl start pbs_sched.service
