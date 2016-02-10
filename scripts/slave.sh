#!/bin/bash
yum -y install wget

# Install torque-mom
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1.2/torque-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque.rpm
wget https://github.com/dizk/torquebuilder/releases/download/v6.0.0.1.2/torque-client-6.0.0.1-1.adaptive.el7.centos.x86_64.rpm -qO torque-client.rpm

yum -y --nogpgcheck localinstall torque.rpm torque-client.rpm

cat > /var/spool/torque/mom_priv/config <<EOF
\$pbsserver      master
\$logevent       255
EOF

# Start pbs-mom
systemctl enable pbs_mom
systemctl start pbs_mom

# Start docquer
cp /home/vagrant/sync/scripts/docquer.service /lib/systemd/system/
systemctl start docquer
