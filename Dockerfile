# Dataplay
#
# VERSION               0.0.1

FROM     ubuntu
MAINTAINER Ben Cartwright Cox "ben@playgen.com"

# make sure the package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-add-repository -y ppa:chris-lea/node.js
RUN apt-get update

RUN apt-get install -y openssh-server screen gcc mysql-server curl git mercurial make binutils bison build-essential wget python-software-properties nodejs
RUN mkdir /var/run/sshd
#RUN screen -dmS SSHD /usr/sbin/sshd -D
#RUN screen -dmS mysql mysqld_safe
RUN echo 'root:dataplay' |chpasswd
RUN mkdir /build/
RUN mkdir /build/redis
RUN wget http://download.redis.io/releases/redis-2.8.9.tar.gz -O /build/redis/redis-2.8.9.tar.gz
RUN tar -C /build/redis -xzf /build/redis/redis-2.8.9.tar.gz
RUN cd /build/redis/redis-2.8.9 && make
RUN screen -dmS redis /build/redis/redis-2.8.9/src/redis-server
RUN curl -s https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar -v -C /usr/local -xz
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /build/DataPlay/
RUN git clone --recursive https://github.com/playgenhub/DataPlay.git /build/DataPlay
RUN sh -c "cd /build/DataPlay; go get; go build"


EXPOSE 22 3000 3306
CMD screen -dmS SSHD /usr/sbin/sshd -D && screen -dmS mysql mysqld_safe && sleep 10 && cat /build/DataPlay/layout.sql | mysql && /build/DataPlay/start.sh