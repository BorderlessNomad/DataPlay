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

def home(l):
	l.client.get("/home/data", headers=headers)
	l.client.get("/chart/toprated", headers=headers)
	l.client.get("/user/activitystream", headers=headers)
	l.client.get("/recentobservations", headers=headers)

class UserBehavior(TaskSet):
	tasks = {home:100}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=500
	max_wait=500
