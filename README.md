Docquer - Homogeneous Torque
================
A simple torque cluster running the work inside a docker container. Comes with a Vagrant cluster if you want to test it.

It uses a simple torque wrapper script in qsubfilters/docker_pbs_wrapper.sh. It just copies the script sent in by qsub and executes in an ubuntu:14.04 container.


_______________________________________________________________________
Build/Run Requirements
----------------------

- Vagrant
- Virtualbox
- Go

Setup
-----
The following command builds a three part virtual Torque cluster with 1 master host and 2 slaves.

    $ NODES=2 vagrant up

`NODES` defines the number of slave nodes that will be created for the cluster.

Access
------
	$ vagrant ssh master
	$ echo 'echo "Hello Docquer"' | qsub

Suspend
-------
	$ NODES=2 vagrant suspend

Take Down
---------
    $ NODES=2 vagrant destroy


Vagrant setup is built upon crcollins/torquecluster
