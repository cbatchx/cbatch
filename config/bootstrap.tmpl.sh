#!/bin/bash

GROUPADD="/usr/sbin/groupadd -g {{.User.GID}} {{.User.Groupname}}"
ADDUSER="/usr/sbin/adduser -d {{.User.HomeDir}} -u {{.User.UID}} -g {{.User.GID}} {{.User.Username}}"
SU="/bin/su"

# µCernVM specifics
# This uses the busybox version of groupadd and adduser. Because the normal version causes parrot to crash.
if [ -f "/UCVM/busybox" ]; 
then
    GROUPADD="/UCVM/busybox addgroup -g {{.User.GID}} {{.User.Groupname}}"
    ADDUSER="/UCVM/busybox adduser -D -h {{.User.HomeDir}} -u {{.User.UID}} -G {{.User.Groupname}} {{.User.Username}}"
    SU="/UCVM/busybox su"
    # Bug with CernVM for docker tmp is	not writable
    /UCVM/busybox chmod 1777 /tmp

fi

# Create a group for the new user.
eval $GROUPADD

# Create the new user.
eval $ADDUSER

# Run the job as the new user.
$SU -m -c "$@" {{.User.Username}}