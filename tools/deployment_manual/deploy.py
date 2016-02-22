#!/usr/bin/python

import getpass
import os
import paramiko
import re
import select
import stat
import sys
import threading
import time
import urllib2

__author__ = 'Mayur Ahir'

# IPs of Servers
REDIS_HOST = '109.231.121.17'
DATABASE_HOST = '109.231.121.123'
CASSANDRA_HOST = '109.231.121.6'

LOADBALANCER_HOST = '109.231.121.141'
POSTGRESQL_HOSTS = ['109.231.121.56', '109.231.121.46']
FRONTEND_HOSTS = ['109.231.121.5']
MASTER_HOSTS = ['109.231.121.10', '109.231.121.71']

# URLs of Deployment scripts
REDIS_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/redis.sh'
PGPOOL_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/pgpool.sh'
CASSANDRA_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/cassandra.sh'

LOADBALANCER_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/loadbalancer/haproxy.sh'
POSTGRESQL_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/db/postgresql.sh'
FRONTEND_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/app/frontend.sh'
MASTER_SCRIPT_URL = 'https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/app/master.sh'

# Other local variables
DIRECTORY = 'scripts'

paramiko.util.log_to_file('deployment.log')
root_command = "whoami\n"
root_command_result = "root"
ssh_pass = getpass.getpass(prompt="Private Key Password? ")

def download_file(directory, url):
    print "DOWNLOADING FILE %s\n" % (url)
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
    print "REPLACE %s with %s in %s\n" % (variable, value, file)
    os.chdir(directory)

    tempFile = file + '.tmp'
    inputFile = open(file)
    outputFile = open(tempFile, 'w')
    fContent = unicode(inputFile.read(), "utf-8")

    regex = re.compile(r'^{}=(.*)$'.format(variable), re.MULTILINE)
    outText = regex.sub(variable + "=\"" + value + "\"", fContent)

    outputFile.write((outText))

    outputFile.close()
    inputFile.close()

    if os.path.isfile(file + ".old"):
        os.remove(file + ".old")
    os.rename(file, file + ".old")
    os.rename(tempFile, file)

def connect_ssh(hostname, username):
    print 'CONNECT SSH, %s@%s\n' % (username, hostname)
    # Create an SSH client
    ssh = paramiko.SSHClient()

    # Make sure that we add the remote server's SSH key automatically
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    # Connect to the client
    ssh.connect(hostname, username=username, password=ssh_pass)

    return ssh

def send_command(ssh, cmd, wait_time, should_print):
    out = ""

    transport = ssh.get_transport()

    channel = transport.open_session()
    channel.exec_command(cmd)

    # Wait for the command to terminate & Print the receive buffer, if necessary
    while should_print and True:
        try:
            if channel.exit_status_ready():
                break
            rl, wl, xl = select.select([channel], [], [], 0.0)
            if len(rl) > 0:
                out = channel.recv(1024)
                if out and out.strip() and not out.isspace():
                    print out
        except KeyboardInterrupt:
            print("Caught CTRL+C")
            channel.close()
            exit(0)

    return out

def send_file(ssh, source_dir, source_file, dest_file, make_executable=True):
    print 'SENDING FILE, %s/%s -> %s\n' % (source_dir, source_file, dest_file)

    # Connect to SFTP client
    sftp = ssh.open_sftp()

    # Transfer local file to remote destination
    sftp.put(source_dir + "/" + source_file, dest_file)

    if make_executable:
        print 'SENDING FILE, chmod +x\n'
        # Make file executable
        sftp.chmod(dest_file, stat.S_IRWXU | stat.S_IRWXG | stat.S_IROTH)

    # Close the SFTP connection
    sftp.close()

def download_send(directory, script_url, host, username, script, dest_path, update_system = False, log_to_file = False):
    download_file(directory, script_url)

    ssh = connect_ssh(host, username)
    send_file(ssh, directory, script, dest_path)

    if update_system:
        if username is 'centos':
            send_command(ssh, "sudo yum update", 1, True)
        elif username is 'ubuntu':
            send_command(ssh, "sudo apt-get update", 1, True)
        else:
            print 'Invalid OS'

    if log_to_file:
        cmd = 'sudo bash ' + dest_path + ' > ' + script + '.log 2>&1 &'
    else:
        cmd = 'sudo bash ' + dest_path
    send_command(ssh, cmd, 1, True)
    ssh.close()
    print

def send(directory, script_url, host, username, script, dest_path, log_to_file = True):
    ssh = connect_ssh(host, username)
    send_file(ssh, directory, script, dest_path)

    send_command(ssh, "sudo apt-get install dos2unix", 1, True)
    send_command(ssh, "dos2unix -k -o " + dest_path, 1, True)

    if log_to_file:
        cmd = 'sudo bash ' + dest_path + ' > ' + script + '.log 2>&1 &'
    else:
        cmd = 'sudo bash ' + dest_path
    send_command(ssh, cmd, 1, True)

    ssh.close()

def task_haproxy(directory):
    download_file(directory, LOADBALANCER_SCRIPT_URL)
    replace_string(directory, 'haproxy.sh', 'REDIS_HOST', REDIS_HOST)

    cmd = threading.Thread(target = send, args = (directory, LOADBALANCER_SCRIPT_URL, LOADBALANCER_HOST, 'ubuntu', 'haproxy.sh', '/home/ubuntu/haproxy.sh'))
    cmd.start()
    cmd.join()

def task_postgresql(directory):
    download_file(directory, POSTGRESQL_SCRIPT_URL)
    replace_string(directory, 'postgresql.sh', 'PGPOOL_API_HOST', DATABASE_HOST)

    for POSTGRESQL_HOST in POSTGRESQL_HOSTS:
        cmd = threading.Thread(target = send, args = (directory, POSTGRESQL_SCRIPT_URL, POSTGRESQL_HOST, 'ubuntu', 'postgresql.sh', '/home/ubuntu/postgresql.sh'))
        cmd.start()
        cmd.join()

def task_frontend(directory):
    download_file(directory, FRONTEND_SCRIPT_URL)
    replace_string(directory, 'frontend.sh', 'LOADBALANCER_HOST', LOADBALANCER_HOST)

    for FRONTEND_HOST in FRONTEND_HOSTS:
        cmd = threading.Thread(target = send, args = (directory, FRONTEND_SCRIPT_URL, FRONTEND_HOST, 'ubuntu', 'frontend.sh', '/home/ubuntu/frontend.sh'))
        cmd.start()
        cmd.join()

def task_master(directory):
    download_file(directory, MASTER_SCRIPT_URL)
    replace_string(directory, 'master.sh', 'DATABASE_HOST', DATABASE_HOST)
    replace_string(directory, 'master.sh', 'REDIS_HOST', REDIS_HOST)
    replace_string(directory, 'master.sh', 'CASSANDRA_HOST', CASSANDRA_HOST)
    replace_string(directory, 'master.sh', 'LOADBALANCER_HOST', LOADBALANCER_HOST)

    for MASTER_HOST in MASTER_HOSTS:
        cmd = threading.Thread(target = send, args = (directory, MASTER_SCRIPT_URL, MASTER_HOST, 'ubuntu', 'master.sh', '/home/ubuntu/master.sh'))
        cmd.start()
        cmd.join()

def main():
    ssh_pass = ''

    directory = os.path.join(os.path.dirname(os.path.abspath(__file__)), DIRECTORY)
    if not os.path.exists(directory):
        os.makedirs(directory)

    # Step 1
    cassandra = threading.Thread(target = download_send, args = (directory, CASSANDRA_SCRIPT_URL, CASSANDRA_HOST, 'ubuntu', 'cassandra.sh', '/home/ubuntu/cassandra.sh', False, True))
    pgpool = threading.Thread(target = download_send, args = (directory, PGPOOL_SCRIPT_URL, DATABASE_HOST, 'centos', 'pgpool.sh', '/home/centos/pgpool.sh', True))
    redis = threading.Thread(target = download_send, args = (directory, REDIS_SCRIPT_URL, REDIS_HOST, 'ubuntu', 'redis.sh', '/home/ubuntu/redis.sh', False, True))

    cassandra.start()
    pgpool.start()
    redis.start()

    cassandra.join()
    pgpool.join()
    redis.join()

    # Step 2
    haproxy = threading.Thread(target = task_haproxy, args = (directory,))
    postgresql = threading.Thread(target = task_postgresql, args = (directory,))

    haproxy.start()
    postgresql.start()

    haproxy.join()
    postgresql.join()

    # Step 3
    frontend = threading.Thread(target = task_frontend, args = (directory,))
    master = threading.Thread(target = task_master, args = (directory,))

    frontend.start()
    master.start()

    frontend.join()
    master.join()

if __name__ == "__main__":
    main()
