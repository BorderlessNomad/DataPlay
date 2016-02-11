#!/usr/bin/python

import getpass
import os
import paramiko
import re
import stat
import sys
import time
import urllib2

__author__ = 'Mayur Ahir'

# IPs of Servers
REDIS_HOST = '109.231.121.5'
DATABASE_HOST = '109.231.121.123'
CASSANDRA_HOST = '109.231.121.56'

LOADBALANCER_HOST = '109.231.121.141'
POSTGRESQL_HOSTS = ['109.231.121.87']
FRONTEND_HOSTS = ['109.231.121.107']
MASTER_HOSTS = ['109.231.121.212']

# URLs of Deployment scripts
REDIS_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/redis.sh'
PGPOOL_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/pgpool.sh'
CASSANDRA_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/cassandra.sh'

LOADBALANCER_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/loadbalancer/haproxy.sh'
POSTGRESQL_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/postgresql.sh'
FRONTEND_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/app/frontend.sh'
MASTER_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/app/master.sh'

def download_file(directory, url):
    print "DOWNLOADING FILE " + url
    os.chdir(directory)

    # file to be written to
    file_name = url.split('/')[-1]

    data = urllib2.urlopen(url).read()

    #open the file for writing
    file = open(file_name, 'wb')

    # read from request while writing to file
    file.write(data)

    file.close()

def replace_string(directory, file, variable, value):
    print "REPLACE %s with %s in %s" % (variable, value, file)
    os.chdir(directory)

    tempFile = file + '.tmp'
    inputFile = open(file)
    outputFile = open(tempFile, 'w')
    fContent = unicode(inputFile.read(), "utf-8")

    regex = re.compile(r'^{}=(.*)$'.format(variable), re.MULTILINE)
    outText = regex.sub(variable + "=\"" + value + "\"", fContent)

    outputFile.write((outText.encode("utf-8")))

    outputFile.close()
    inputFile.close()

    if os.path.isfile(file + ".old"):
        os.remove(file + ".old")
    os.rename(file, file + ".old")
    os.rename(tempFile, file)

def send_string_and_wait(ssh, command, wait_time, should_print):
    # Create a raw shell
    shell = ssh.invoke_shell()

    # Send the su command
    shell.send(command)

    # Wait a bit, if necessary
    time.sleep(wait_time)

    # Flush the receive buffer
    receive_buffer = shell.recv(1024)

    # Print the receive buffer, if necessary
    if should_print:
        print receive_buffer

def send_file(hostname, username, source_dir, source_file, dest_file, make_executable=True):
    print 'SENDING FILE, %s/%s -> %s = %s@%s' % (source_dir, source_file, dest_file, username, hostname)
    # Create an SSH client
    ssh = paramiko.SSHClient()

    # Make sure that we add the remote server's SSH key automatically
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    # Connect to the client
    ssh.connect(hostname, username=username, password=ssh_pass)

    # Create a raw shell
    shell = ssh.invoke_shell()

    # Send the su command
    # send_string_and_wait(ssh, "sudo su\n", 1, True)

    # Send command on root
    # send_string_and_wait(ssh, root_command, 1, True)

    # Connect to SFTP client
    sftp = ssh.open_sftp()

    # Transfer local file to remote destination
    sftp.put(source_dir + "/" + source_file, dest_file)

    if make_executable:
        print 'SENDING FILE, chmod +x'
        # Make file executable
        sftp.chmod(dest_file, stat.S_IRWXU | stat.S_IRWXG | stat.S_IROTH)

    # Close the SFTP connection
    sftp.close()

    # Close the SSH connection
    ssh.close()

def main():
    ssh_pass = ''

    directory = os.path.join(os.path.dirname(os.path.abspath(__file__)), DIRECTORY)
    if not os.path.exists(directory):
        os.makedirs(directory)

    # Cassandra
    download_file(directory, CASSANDRA_SCRIPT_URL)
    send_file(CASSANDRA_HOST, 'ubuntu', directory, 'cassandra.sh', '/home/ubuntu/cassandra.sh')
    print

    # PgPool-II
    download_file(directory, PGPOOL_SCRIPT_URL)
    send_file(DATABASE_HOST, 'centos', directory, 'pgpool.sh', '/home/centos/pgpool.sh')
    print

    # Redis
    download_file(directory, REDIS_SCRIPT_URL)
    send_file(REDIS_HOST, 'ubuntu', directory, 'redis.sh', '/home/ubuntu/redis.sh')
    print

    # HAProxy
    download_file(directory, LOADBALANCER_SCRIPT_URL)
    replace_string(directory, 'haproxy.sh', 'REDIS_HOST', REDIS_HOST)
    send_file(LOADBALANCER_HOST, 'ubuntu', directory, 'haproxy.sh', '/home/ubuntu/haproxy.sh')
    print

    # PostgreSQL
    download_file(directory, POSTGRESQL_SCRIPT_URL)
    replace_string(directory, 'postgresql.sh', 'PGPOOL_API_HOST', DATABASE_HOST)
    # send_file(REDIS_HOST, 'ubuntu', directory, 'redis.sh', '/home/ubuntu/redis.sh')
    print

    # Frontend
    download_file(directory, FRONTEND_SCRIPT_URL)
    replace_string(directory, 'frontend.sh', 'LOADBALANCER_HOST', LOADBALANCER_HOST)
    # send_file(REDIS_HOST, 'ubuntu', directory, 'redis.sh', '/home/ubuntu/redis.sh')
    print

    # Master
    download_file(directory, MASTER_SCRIPT_URL)
    replace_string(directory, 'master.sh', 'DATABASE_HOST', DATABASE_HOST)
    replace_string(directory, 'master.sh', 'REDIS_HOST', REDIS_HOST)
    replace_string(directory, 'master.sh', 'CASSANDRA_HOST', CASSANDRA_HOST)
    replace_string(directory, 'master.sh', 'LOADBALANCER_HOST', LOADBALANCER_HOST)
    # send_file(REDIS_HOST, 'ubuntu', directory, 'redis.sh', '/home/ubuntu/redis.sh')
    print

# Other local variables
DIRECTORY = 'scripts'

paramiko.util.log_to_file('demo.log')
root_command = "whoami\n"
root_command_result = "root"
ssh_pass = getpass.getpass(prompt="Private Key Password? ")

if __name__ == "__main__":
    main()
