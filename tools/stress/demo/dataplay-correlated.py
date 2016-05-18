from locust import HttpLocust, TaskSet

session = ""
headers = {'X-API-SESSION': ''}

def login(l):
	global headers
	global session

	if not session:
		r = l.client.post("/api/login", json={'username': 'mayur', 'password': '123456'})
		response = r.json()

		session = response["session"]
		headers = {'X-API-SESSION': session}

def correlated(l):
	l.client.get("/api/correlated/births", headers=headers)
	l.client.get("/api/correlated/job-density", headers=headers)
	l.client.get("/api/correlated/deaths", headers=headers)
	l.client.get("/api/correlated/gold-prices", headers=headers)
	l.client.get("/api/correlated/jobs", headers=headers)

class UserBehavior(TaskSet):
	tasks = {correlated}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=1000
	max_wait=1000
