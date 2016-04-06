############################################################
# Dockerfile to build testAppGru container images
# Based on Ubuntu
############################################################

# Set the base image to Ubuntu
FROM ubuntu

# File Author / Maintainer
MAINTAINER Luca Florio

################## BEGIN INSTALLATION ######################
# Usage: ADD [source directory or URL] [destination directory]
ADD testAppGru /usr/bin/testAppGru

##################### INSTALLATION END #####################

# Set default container command
ENTRYPOINT ["testAppGru"]