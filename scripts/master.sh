#!/bin/bash
apt-get install -y torque-server torque-scheduler  torque-client

qterm
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


for i in `seq 1 $1`; do echo  "slave$i np=1" >> /var/spool/torque/server_priv/nodes; done

qterm
pbs_server
