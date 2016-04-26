#!/bin/bash

# Make config
cat > /etc/cvmfs/default.local <<EOF
CVMFS_REPOSITORIES=atlas.cern.ch,atlas-condb.cern.ch,grid.cern.ch,cernvm-prod.cern.ch
CVMFS_HTTP_PROXY="DIRECT"
EOF

# Setup cvmfs
service autofs start
service autofs reload
service autofs restart
cvmfs_config setup

cat > /etc/auto.master <<EOF
/cvmfs /etc/auto.cvmfs
EOF

# Test autofs
cvmfs_config probe
service autofs stop

/usr/sbin/automount -v -f