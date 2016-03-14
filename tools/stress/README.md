DataPlay Load Testing Tool
===

Installation
---
 1. Install latest version of Python 2.x
 2. Install **pip** package management system
 3. Install Locust.io using command `pip install locustio`


Usage
---
Run `$ locust -f dataplay.py --host=URL` in your Bash or CMD prompt.
 
 **URL** can be of,

 - Load Balancer e.g. `$ locust -f dataplay.py --host=http://109.231.121.51`
 - Master e.g. `$ locust -f dataplay.py --host=http://109.231.121.61:3000`

Then, open your browser & navigate to `http://127.0.0.1:8089/`

To Simulate load in background,

```
nohup locust -f dataplay.py --host=http://109.231.121.52/api --no-web --print-stats --clients=100 --hatch-rate=5 --loglevel=INFO --logfile=dataplay.debug.log > dataplay.web.log 2>&1&
```