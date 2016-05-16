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

def correlated_generate(l):
	l.client.get("/correlated/births/true", headers=headers)
	l.client.get("/correlated/job-density/true", headers=headers)
	l.client.get("/correlated/deaths/true", headers=headers)
	l.client.get("/correlated/gold-prices/true", headers=headers)
	l.client.get("/correlated/jobs/true", headers=headers)

class UserBehavior(TaskSet):
	tasks = {correlated_generate:100}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=500
	max_wait=500
