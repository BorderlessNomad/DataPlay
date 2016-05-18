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

def related(l):
	l.client.get("/api/related/births", headers=headers)
	l.client.get("/api/related/job-density", headers=headers)
	l.client.get("/api/related/deaths", headers=headers)
	l.client.get("/api/related/gold-prices", headers=headers)
	l.client.get("/api/related/jobs", headers=headers)

class UserBehavior(TaskSet):
	tasks = {related}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=1000
	max_wait=1000
