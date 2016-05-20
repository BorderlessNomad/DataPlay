'''
Created on Mar 1, 2016

@author: ahmeda
'''
import math, json
import random
from time import sleep
from requests.auth import HTTPBasicAuth

import random
from requests import Request, Session

import multiprocessing as mp
import time, requests, math, traceback
#from __main__ import traceback
import numpy as np

random.seed(0)


def nextTime(rateParameter):
    return -math.log(1.0 - random.random()) / rateParameter


TimeOfReq = []
for i in range(5, 161, 5):
    for j in range(10):
        summer = 0
        for k in range(1, i):
            const = nextTime(i)
            TimeOfReq += [const]
            summer += const
        if summer < 1:
            TimeOfReq[-1] = 1 - summer + TimeOfReq[-1]
# elif summer>1:
#                 x[-1]=1-summer-x[-1]

# HOST='http://p03.ds.cs.umu.se:8092/gw'
# HOST='http://p03.ds.cs.umu.se:8092/PHP/index.html'
# LINK='index.php/Adolf_Hitler'

CONSTANT = 25
TIME_OF_RUNNING = 30

DROP_REQUESTS_AFTER = 300
SLEEP = 1

auth = json.dumps({'username': 'mayur', 'password': '123456'})
RAMP = range(1, 200, 1)
RAMP += [15] * 30
RAMP *= 5
CONSTANT_WORKLOAD = [CONSTANT] * TIME_OF_RUNNING

NUMBER_OF_REQUESTS = CONSTANT_WORKLOAD
# Links=open("links.out",'r')
# links=Links.readlines()
# links=[l.strip() for l in links]
# print links

HOST = "http://109.231.121.115"
login = HOST + "/api/login"
r = requests.post(url=login, data=auth)
XSessionID = r.content.split("\"")[3]

url = []
url += [HOST + "/api/correlated/deaths/true"]
url += [HOST + "/api/correlated/births/true"]
url += [HOST + "/api/correlated/job-density/true"]
url += [HOST + "/api/correlated/gold-prices/true"]

fn = './CactosLog.out'

#
def worker(j, q, u):
    '''stupidly simulates long running process'''
    #     for weights in ['uniform', 'distance']:
    global links
    t1 = time.time()
    s = Session()
    req = Request('GET', u, headers={'X-API-SESSION': XSessionID})
    prepped = s.prepare_request(req)

    resp = s.send(prepped)
    t2 = time.time()

    RT = t2 - t1
    out = u, str(RT) + "  " + str(j) + "   " + str(t1) + "   ", r.status_code, " ", t1, t2
    q.put(out)
    print "out", out
    return out


def listener(q):
    '''listens for messages on the q, writes to file. '''

    f = open(fn, 'a')
    while 1:
        m = q.get()
        #    print "m", m
        if m == 'kill':
            f.write('killed')
            break
        f.write(str(m) + '\n')
        f.flush()
    f.close()


def main():
    # must use Manager queue here, or will not work
    manager = mp.Manager()
    q = manager.Queue()
    pool = mp.Pool(1000)

    # put listener to work first
    watcher = pool.apply_async(listener, (q,))
    random.seed(0)

    jobs = []
    # fire off workers
    for i in NUMBER_OF_REQUESTS:
        time.sleep(SLEEP)
        print i
        for j in range(1, i + 1):
            job = pool.apply_async(worker, (i, q, random.choice(url)))
            jobs.append(job)

    # collect results from the workers through the pool result queue
    for job in jobs:
        try:
            job.get()
        #        print xxx
        except:
            print "There was an error"
            print traceback.print_exc()
            continue

    # now we are done, kill the listener
    q.put('kill')
    pool.close()


if __name__ == "__main__":
    main()
