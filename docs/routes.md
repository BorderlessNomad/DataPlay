HTTP Routes
===

### /noauth/login.json
**POST - No auth needed**

**Handled in main:HandleLogin:main.go**

Expects a HTTP form input that has `username` and `password`.
Will then set the login cookie and redirect the user to the root of the domain `/`

```html

        <form class="form-signin" action="/noauth/login.json" method="POST" >
          <input name="username" type="text" class="form-control" placeholder="Email address" required autofocus>
          <input name="password" type="password" class="form-control" placeholder="Password" required>
          <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
        </form>

```

### /noauth/register.json
**POST - No auth needed**

**Handled in main:HandleRegister:main.go**

Expects a HTTP form input that has `username` and `password`.
In the case that it the user registering already exists, then the service will respond back with

`That username is already registered.`

eles it will 302 redirect to the homepage with the new user logged in, ready to go.

### /api/user
**GET - Auth Cookie Checked**

**Handled in api:CheckAuth:api/datalisting.go**

This is used to get the user infomation, it can also be used as a way to check if the user is still logged into the server.

The responce is a json encoded payload that follows:

```json

	{
	    "Username": "ben@playgen.com",
	    "UserID": 2
	}

```

### /api/visited
**GET - Auth Cookie Checked**

**Handled in api:GetLastVisited:api/tracking.go**

This is used to get the last datasets that the authenticated user visited, the top 5 are given back in the form of a JSON array.

an example payload is as follows:

```json

	[
		"hips",
		"Hip Fractures",
		"false"
	],

```

`hips` is the GUID of the dataset that can be used on the URL slug. `Hip Fractures` is the friendly name of the dataset that should be displayed to the user and the final string `false` is if the dataset contains location data, if it is `true` then you need to offer the user a map view of the data.


```json

	[
	    [
	        "2e933891cf7c8e5db3704632c6a0b72cde17e488a90c4b68c43ececcda5",
	        "DCMS Government Procurement Card Spend أƒآ¢أ¢â€ڑآ¬أ¢â‚¬إ“Transactions over أƒâ€ڑأ‚آ£500, 1 April 2011 - October 2011",
	        "false"
	    ],
	    [
	        "hips",
	        "Hip Fractures",
	        "false"
	    ],
	    [
	        "weather_uk",
	        "UK Weather",
	        "true"
	    ],
	    [
	        "c47dea484967a26ee56bd326fc1d361b70f84c70ff13b3b8c17771fc033",
	        "Contracts and Tenders from Bristol City Council",
	        "false"
	    ],
	    [
	        "547edfc17a8e430c292e6d3da78b14bb1e8b640dd82ed9c89b7f55d8c8c",
	        "April 2013 Spend over أƒâ€ڑأ‚آ£25k",
	        "false"
	    ]
	]

```

### /api/search/:s
**GET - Auth Cookie Checked**

**Handled in api:SearchForData:api/datalisting.go**

This call is used to search the index for datasets that match words.

A example query such as `api/search/bris` will return the following array of json payloads:

```json

	[
	    {
	        "Title": "Contracts and Tenders from Bristol City Council",
	        "GUID": "066b4f97fa16d12af8be2b1be67e8d2d363b2807ff13b3b8c17771fc033",
	        "LocationData": "false"
	    },
	    {
	        "Title": "Contracts and Tenders from Bristol City Council",
	        "GUID": "2ece89a545d7e610b3878ce2ba7ba4f57f3e96c2ff13b3b8c17771fc033",
	        "LocationData": "false"
	    }
	]

```

### /api/getinfo/:id
**GET - Auth Cookie Checked**

**Handled in api:GetEntry:api/datalisting.go**

### /api/getimportstatus/:id
**GET - Auth Cookie Checked**

**Handled in api:CheckImportStatus:api/datalisting.go**

### /api/getdata/:id
**GET - Auth Cookie Checked**

**Handled in api:DumpTable:api/datalisting.go**

This function dumps the **whole** table as a array of json objects, Please make sure this is what you want to do before doing it!

Here is a example output (with 99% of the data ommitted):

```json

	[
	    {
	        "date": "1950-01-01",
	        "price": "34.73"
	    },
	    {
	        "date": "1950-02-01",
	        "price": "34.73"
	    },
	    {
	        "date": "1950-03-01",
	        "price": "34.73"
	    }
	]

```

Please note even if the output module will automatically cast *everything* that comes out of the database as a string.

### /api/getdata/:id/:top/:bot
**GET - Auth Cookie Checked**

**Handled in api:DumpTable:api/datalisting.go**

This works similar to `getdata` above however has a way to load parts of the data at a time.

the `:top` part is what row you want to start reading and the `:bot` part is how many rows you want to read for.

For example the call: `/api/getdata/gold/0/1` will return

```json

	[
	    {
	        "date": "1950-01-01",
	        "price": "34.73"
	    }
	]

```

and `/api/getdata/gold/1/1` will return the next entry in the dataset:


```json

	[
	    {
	        "date": "1950-02-01",
	        "price": "34.73"
	    }
	]

```

### /api/getdata/:id/:x/:startx/:endx
**GET - Auth Cookie Checked**

**Handled in api:DumpTableRange:api\datalisting.go**

This allows you to select a part of a dataset to start reading from where `:x` is the colloum for a **numeric** value.

For example if I want to select years 1995 to 2000 in the dataset `GDP` I would do `/api/getdata/gdp/year/1995/2000` and it would return:

```json

	[
	    {
	        "GDP": "1005050",
	        "change": "3",
	        "gdpindex": "71.7",
	        "year": "1995"
	    },
	    {
	        "GDP": "1036340",
	        "change": "3",
	        "gdpindex": "73.9",
	        "year": "1996"
	    },
	    {
	        "GDP": "1076350",
	        "change": "4",
	        "gdpindex": "76.8",
	        "year": "1997"
	    },
	    {
	        "GDP": "1114180",
	        "change": "4",
	        "gdpindex": "79.5",
	        "year": "1998"
	    },
	    {
	        "GDP": "1149460",
	        "change": "3",
	        "gdpindex": "82",
	        "year": "1999"
	    },
	    {
	        "GDP": "1198150",
	        "change": "4",
	        "gdpindex": "85.5",
	        "year": "2000"
	    }
	]

```


### /api/getdatagrouped/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:DumpTableGrouped:api\datalisting.go**

### /api/getdatapred/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:DumpTablePrediction:api\datalisting.go**

### /api/getcsvdata/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:GetCSV:api/datalisting.go**

### /api/getreduceddata/:id
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

### /api/getreduceddata/:id/:persent
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

### /api/getreduceddata/:id/:persent/:min
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

### "/api/setdefaults/:id
**GET - Auth Cookie Checked**

**Handled in api:SetDefaults:api/**

### /api/getdefaults/:id
**GET - Auth Cookie Checked**

**Handled in api:GetDefaults:api/**

### /api/identifydata/:id
**GET - Auth Cookie Checked**

**Handled in api:IdentifyTable:api/**

### /api/findmatches/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:AttemptToFindMatches:api/**

### /api/classifydata/:table/:col
**GET - Auth Cookie Checked**

**Handled in api:SuggestColType:api/**

### /api/stringmatch/:word
**GET - Auth Cookie Checked**

**Handled in api:FindStringMatches:api/**

### /api/stringmatch/:word/:x
**GET - Auth Cookie Checked**

**Handled in api:FindStringMatches:api/**

### /api/relatedstrings/:guid
**GET - Auth Cookie Checked**

**Handled in api:GetRelatedDatasetByStrings:api/**
