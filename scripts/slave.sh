#!/bin/bash
yum -y install epel-release

mkdir -p /var/lib/torque/mom_priv
cat > /var/lib/torque/mom_priv/config <<EOF
\$pbsserver      master
\$logevent       255
\$usecp *:/home  /mnt/nfs/home/
EOF

yum -y install torque-mom

# Setup NFS
mkdir -p /mnt/nfs/home
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server
mount -t nfs 192.168.1.100:/home /mnt/nfs/home/
echo "192.168.1.100:/home    /mnt/nfs/home   nfs defaults 0 0" >> /etc/fstab


# Start pbs-mom
systemctl enable pbs_mom
systemctl start pbs_mom

# Start docquer
cp /vagrant/scripts/docquer.service /lib/systemd/system/
# systemctl start docquer
