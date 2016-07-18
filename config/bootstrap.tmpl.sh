#!/bin/bash

# Create a group for the new user.
/usr/sbin/groupadd -g {{.User.GID}} {{.User.Groupname}}

# Create the new user.
/usr/sbin/adduser -d {{.User.HomeDir}} -u {{.User.UID}} -g {{.User.GID}} {{.User.Username}}

# Run the job as the new user.
su -c 'exec $@' {{.User.Username}} -- "$@"