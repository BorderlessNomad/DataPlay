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

def activities(l):
	l.client.get("/political/popular", headers=headers)
	l.client.get("/political/keywords", headers=headers)
	l.client.get("/political/mediapulse", headers=headers)
	l.client.get("/political/regions", headers=headers)

class UserBehavior(TaskSet):
	tasks = {activities:100}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=500
	max_wait=500
