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
	        "DCMS Government Procurement Card Spend Ã¢â‚¬â€œTransactions over Ã‚Â£500, 1 April 2011 - October 2011",
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
	        "April 2013 Spend over Ã‚Â£25k",
	        "false"
	    ]
	]

```

### /api/search/:s
**GET - Auth Cookie Checked**

**Handled in api:SearchForData:api/datalisting.go**

### /api/getinfo/:id
**GET - Auth Cookie Checked**

**Handled in api:GetEntry:api/datalisting.go**

### /api/getimportstatus/:id
**GET - Auth Cookie Checked**

**Handled in api:CheckImportStatus:api/datalisting.go**

### /api/getdata/:id
**GET - Auth Cookie Checked**

**Handled in api:DumpTable:api/datalisting.go**

### /api/getdata/:id/:top/:bot
**GET - Auth Cookie Checked**

**Handled in api:DumpTable:api/datalisting.go**

### /api/getdata/:id/:x/:startx/:endx
**GET - Auth Cookie Checked**

**Handled in api:DumpTableRange:api\datalisting.go**

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
