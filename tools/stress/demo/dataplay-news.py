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

def news(l):
	l.client.get("/news/search/births", headers=headers)
	l.client.get("/news/search/death", headers=headers)
	l.client.get("/news/search/population", headers=headers)
	l.client.get("/news/search/health", headers=headers)
	l.client.get("/news/search/newham", headers=headers)
	l.client.get("/news/search/westminster", headers=headers)

class UserBehavior(TaskSet):
	tasks = {news:100}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=500
	max_wait=500
