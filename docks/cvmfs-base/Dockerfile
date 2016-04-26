FROM ubuntu:15.04

# Prevent debconf errors
ENV DEBIAN_FRONTEND=noninteractive

# Install wget
RUN apt-get update && apt-get install -y wget

# Install cernvmfs repository.   
RUN  wget -q https://ecsft.cern.ch/dist/cvmfs/cvmfs-release/cvmfs-release-latest_all.deb \
     && dpkg -i cvmfs-release-latest_all.deb \
     && rm -f cvmfs-release-latest_all.deb \
     && apt-get update
     
# Install cernvmfs
RUN apt-get install -y \
    cvmfs \
    cvmfs-config-default \
    uuid-runtime 


COPY entry.sh /
ENTRYPOINT ["/entry.sh"]