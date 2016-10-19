From ubuntu:16.04
# Docker file for runnel

# Install required system packages
RUN apt-get -y update
RUN apt-get install -y redis-server iputils-ping

#### Install your project cli dependencies here


# New user for runnel
RUN useradd -ms /bin/bash runnel

# Copy runnel server
COPY bin/linux_64/runnel_server /home/runnel/runnel_server

# Change user to runnel
RUN chown -R runnel:runnel /home/runnel
USER runnel

# Run runnel and redis
EXPOSE 9090
CMD redis-server /etc/redis/redis.conf && /home/runnel/runnel_server -p 0.0.0.0:9090
