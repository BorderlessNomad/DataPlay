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
	l.client.get("/related/births", headers=headers)
	l.client.get("/related/job-density", headers=headers)
	l.client.get("/related/deaths", headers=headers)
	l.client.get("/related/gold-prices", headers=headers)
	l.client.get("/related/jobs", headers=headers)

class UserBehavior(TaskSet):
	tasks = {related:100}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=500
	max_wait=500
