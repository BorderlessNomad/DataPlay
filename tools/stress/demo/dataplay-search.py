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

def search(l):
	l.client.get("/api/search/birth", headers=headers)
	l.client.get("/api/search/death", headers=headers)
	l.client.get("/api/search/population", headers=headers)
	l.client.get("/api/search/health", headers=headers)
	l.client.get("/api/search/newham", headers=headers)
	l.client.get("/api/search/westminster", headers=headers)

class UserBehavior(TaskSet):
	tasks = {search}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=1000
	max_wait=1000
