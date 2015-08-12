from locust import HttpLocust, TaskSet

global headers
headers = {'X-API-SESSION': 'ubVBoLmwvBQoM5GXrb2lKuQWciAg6U0X1RBBBNz7vwv0bhd7RX4JC4n39ztBb2Kl'}

def home(l):
	l.client.get("/home/data", headers=headers)
	l.client.get("/chart/toprated", headers=headers)
	l.client.get("/user/activitystream", headers=headers)
	l.client.get("/recentobservations", headers=headers)

def search(l):
	l.client.get("/search/birth", headers=headers)
	l.client.get("/search/death", headers=headers)
	l.client.get("/search/population", headers=headers)
	l.client.get("/search/health", headers=headers)
	l.client.get("/search/newham", headers=headers)
	l.client.get("/search/westminster", headers=headers)

def news(l):
	l.client.get("/news/search/births", headers=headers)
	l.client.get("/news/search/death", headers=headers)
	l.client.get("/news/search/population", headers=headers)
	l.client.get("/news/search/health", headers=headers)
	l.client.get("/news/search/newham", headers=headers)
	l.client.get("/news/search/westminster", headers=headers)

def related(l):
	l.client.get("/related/births", headers=headers)
	l.client.get("/related/job-density", headers=headers)
	l.client.get("/related/deaths", headers=headers)
	l.client.get("/related/gold-prices", headers=headers)
	l.client.get("/related/jobs", headers=headers)

def correlated(l):
	l.client.get("/correlated/births", headers=headers)
	l.client.get("/correlated/job-density", headers=headers)
	l.client.get("/correlated/deaths", headers=headers)
	l.client.get("/correlated/gold-prices", headers=headers)
	l.client.get("/correlated/jobs", headers=headers)

def correlated_generate(l):
	l.client.get("/correlated/births/true", headers=headers)
	l.client.get("/correlated/job-density/true", headers=headers)
	l.client.get("/correlated/deaths/true", headers=headers)
	l.client.get("/correlated/gold-prices/true", headers=headers)
	l.client.get("/correlated/jobs/true", headers=headers)

def activities(l):
	l.client.get("/political/popular", headers=headers)
	l.client.get("/political/keywords", headers=headers)
	l.client.get("/political/mediapulse", headers=headers)
	l.client.get("/political/regions", headers=headers)

class UserBehavior(TaskSet):
	tasks = {home:100, search:50, news:50, related:25, correlated:25, correlated_generate: 5, activities:10}

class WebsiteUser(HttpLocust):
	task_set = UserBehavior
	min_wait=5000
	max_wait=9000
