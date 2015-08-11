from locust import HttpLocust, TaskSet

def login(l):
    l.client.post("/login", {"username":"mayur", "password":"123456"})

def home(l):
    l.client.get("/home/data")
    l.client.get("/chart/toprated")
    l.client.get("/user/activitystream")
    l.client.get("/recentobservations")

def search(l):
    l.client.get("/search/birth")
    l.client.get("/search/death")
    l.client.get("/search/population")
    l.client.get("/search/health")
    l.client.get("/search/newham")
    l.client.get("/search/westminster")

def news(l):
	l.client.get("/news/search/births")
	l.client.get("/news/search/death")
	l.client.get("/news/search/population")
	l.client.get("/news/search/health")
	l.client.get("/news/search/newham")
	l.client.get("/news/search/westminster")

def related(l):
	l.client.get("/related/births")
	l.client.get("/related/mental-health-problems")
	l.client.get("/related/deaths")
	l.client.get("/related/gold-prices")
	l.client.get("/related/jobs")

def correlated(l):
	l.client.get("/correlated/births")
	l.client.get("/correlated/mental-health-problems")
	l.client.get("/correlated/deaths")
	l.client.get("/correlated/gold-prices")
	l.client.get("/correlated/jobs")

def activities(l):
	l.client.get("/political/popular")
	l.client.get("/political/keywords")
	l.client.get("/political/mediapulse")
	l.client.get("/political/regions")

class UserBehavior(TaskSet):
    tasks = {home:1, search:1, news:1, related:1, correlated:1, activities:1}

    def on_start(self):
        login(self)

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait=5000
    max_wait=9000
