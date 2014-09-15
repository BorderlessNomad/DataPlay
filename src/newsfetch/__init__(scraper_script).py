import logging, json, importio, latch

client = importio.importio(user_id="cf592fba-bd1f-4128-8e98-e729c2bb7dec", api_key="aledxqRLOCLFo9O7cYeeC58aotifmZbL2C57Mg1zicz6ZLVSY94xttvI9AjeV1Fw9DpBg2y/cbrNZXM23yiWBg==", host="https://query.import.io")
client.connect()
queryLatch = latch.latch(13441)
dataRows = []
d = ''

def callback(query, message):
  global dataRows
  global d
  
  if message["type"] == "DISCONNECT":
    print "Query in progress when library disconnected"
    print json.dumps(message["data"], indent = 4)

  if message["type"] == "MESSAGE":
    
    if "errorType" in message["data"]:
      print "Got an error!" 
      print json.dumps(message["data"], indent = 4)
    else:
      print "Got data!"
      print json.dumps(message["data"], indent = 4)
      dataRows.extend(message["data"]["results"])
      d = message["data"]["results"]
      for i in d:
        with open('urls.txt', 'a') as f:
          f.write(i["url"] + ',\n')
 
  if query.finished(): queryLatch.countdown()
  
url = ''

for year in range(5):
  for month in range(1,13):
    for step in range(224):
      if year == 4 and month >= 9:
        break
      url = 'http://www.bbc.co.uk/search/news/?page=' + str(step+1) + '&q=that&text=on&start_day=01&start_month=' + ('{:0>2d}'.format(month)) + '&start_year=201' + str(year) + '&end_day=20&end_month=08&end_year=2014&sort=reversedate&dir=fd&news=' + str(step*20+1) + '&news_av=1'
      client.query({"connectorGuids": ["0db06af4-dd7a-42f2-9672-b8f62e0f98ac"],"input": {"webpage/url": url}}, callback)

print "Queries dispatched, now waiting for results"
queryLatch.await()
print "Latch has completed, all results returned"
client.disconnect()
print "All data received:"
print json.dumps(dataRows, indent = 4)
      

