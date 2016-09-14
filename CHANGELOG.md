0.0.6 - 14.09.2016
### Added
    - Additional logging in the docker driver.
    - Set tighter permissions on boostrap file and set it to executable.
    - Use docker.NewVersionedClientFromEnv("") to get the Docker client env. To enable to lock down of docker API versions.

0.0.5 - 19.07.2016
### Added
    - Support for custom boostrapping scripts
### Removed
    - Does not create a modified image with the new user. And use the new bootstrap script for this.

0.0.4 - 25.04.2016 
### Added
    - Can import from source http
    - Rewrote how users are inserted, now commits a new container to do it.

0.0.3 - 06.04.2016
### Added
    - Takes a snapshot of the environment and copies them to the container.

0.0.2 - 09.03.2016
### Added
    - Rewrite to be a job starter.
### Removed
    - No longer a daemon.

0.0.1 - 14.02.2016
### Added
    - Basic functionality, can run jobs inside containers
