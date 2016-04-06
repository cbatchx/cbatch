#!/bin/bash


echo "Downloading torque rpms..."
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1/torque-client-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-client.rpm

rpm -i torque.rpm torque-client.rpm
rm torque*.rpm

# Setup NFS
mkdir -p /mnt/nfs/home
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server
mount -t nfs 192.168.1.100:/home /mnt/nfs/home/
echo "192.168.1.100:/home    /mnt/nfs/home   nfs defaults 0 0" >> /etc/fstab

# Configure torque MOM.
# Note that job_starter is poiting at /vagrant/cbatch. This is for
# development purposes.
cat > /var/spool/torque/mom_priv/config <<EOF
\$pbsserver      master
\$logevent       255
\$usecp *:/home  /mnt/nfs/home/
\$job_starter /vagrant/cbatch
\$job_starter_run_privileged true
EOF

# Start pbs-mom
sleep 2
systemctl enable pbs_mom
systemctl stop pbs_mom
systemctl start pbs_mom

# Copy in cbatch config
cp /vagrant/config/config.toml /etc/cbatch.toml
