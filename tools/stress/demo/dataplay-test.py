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

def correlated_generate(l):
	l.client.get("/api/correlated/births/true", headers=headers)
	l.client.get("/api/correlated/job-density/true", headers=headers)
	l.client.get("/api/correlated/deaths/true", headers=headers)
	l.client.get("/api/correlated/gold-prices/true", headers=headers)
	l.client.get("/api/correlated/jobs/true", headers=headers)

def search(l):
	l.client.get("/api/search/birth", headers=headers)
	l.client.get("/api/search/death", headers=headers)
	l.client.get("/api/search/population", headers=headers)
	l.client.get("/api/search/health", headers=headers)
	l.client.get("/api/search/newham", headers=headers)
	l.client.get("/api/search/westminster", headers=headers)
	
def related(l):
	l.client.get("/api/related/births", headers=headers)
	l.client.get("/api/related/job-density", headers=headers)
	l.client.get("/api/related/deaths", headers=headers)
	l.client.get("/api/related/gold-prices", headers=headers)
	l.client.get("/api/related/jobs", headers=headers)

class UserBehavior(TaskSet):
	tasks = {correlated:3, search:3, related:3, correlated_generate:1}

	def on_start(self):
		login(self)

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=2000
	max_wait=2000
