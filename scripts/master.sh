#!/bin/bash

# Installing Influxdb
cat <<EOF | sudo tee /etc/yum.repos.d/influxdb.repo
[influxdb]
name = InfluxDB Repository - RHEL \$releasever
baseurl = https://repos.influxdata.com/rhel/\$releasever/\$basearch/stable
enabled = 1
gpgcheck = 1
gpgkey = https://repos.influxdata.com/influxdb.key
EOF

sudo yum -y install influxdb

echo "Downloading torque rpms..."
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-scheduler-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-scheduler.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-server-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-server.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-client-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-client.rpm

rpm -i torque.rpm torque-server.rpm torque-scheduler.rpm torque-client.rpm
rm torque*.rpm

mkdir -p /var/spool/torque/checkpoint/

systemctl enable trqauthd
systemctl start trqauthd


# Configure pbs_server
echo "Configure pbs_server"
yes | pbs_server -t create
sleep 5
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
echo "Creating nodes file"
: > /var/spool/torque/server_priv/nodes
for i in `seq 1 $1`; do echo  "slave$i np=5" >> /var/spool/torque/server_priv/nodes; done

# Set server_name to HOSTNAME
echo $HOSTNAME > /var/spool/torque/server_name

echo "Starting pbs_server"

# Enable and start services
systemctl enable pbs_server
systemctl start pbs_server
systemctl enable pbs_sched
systemctl start pbs_sched

# Configure NFS
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server

# This is only a test setup, this is not secure
echo "/home *(rw,sync,no_root_squash,no_all_squash)" > /etc/exports
systemctl restart nfs-server

# Start influxdb
systemctl start  influxdb.service

# Wait for influxdb
sleep 3

# Insert a test user
influx -execute "create user test with password 'test'"
influx -execute "grant all privileges to test"
