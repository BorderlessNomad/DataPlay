# Dataplay
#
# VERSION               0.0.1

FROM     ubuntu
MAINTAINER Ben Cartwright Cox "ben@playgen.com"

# make sure the package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get update

RUN apt-get install -y openssh-server screen gcc mysql-server curl git mercurial make binutils bison build-essential wget
RUN mkdir /var/run/sshd 
RUN screen -dmS SSHD /usr/sbin/sshd -D
RUN echo 'root:dataplay' |chpasswd
RUN mkdir /build/redis
RUN wget http://download.redis.io/releases/redis-2.8.9.tar.gz -O /build/redis/
RUN tar -C /build/redis -xzf redis-2.8.9.tar.gz
RUN sh -c "cd redis-2.8.9; make"
RUN screen -dmS redis /build/redis/redis-2.8.9/src/redis-server
RUN curl -s https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar -v -C /usr/local -xz
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1
ENV PATH /usr/local/go/bin:$PATH
RUN git clone https://github.com/playgenhub/DataPlay.git /build/
RUN sh -c "cd /build/DataPlay; go get; go build"
RUN cat /build/DataPlay/layout.sql | mysql

EXPOSE 22 3000 3306
CMD /build/DataPlay/DataPlay