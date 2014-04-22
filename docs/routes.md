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

This call groups all the rows that share the same string in `:x` and then adds up all of their `:y` to be one.

For example, If I wanted to find out how many follows I have observed in each CCode in the tweets dataset, I can do:

`/api/getdatagrouped/tweets/CCode/followers`

and I will get back:

```json

	[
	    {
	        "CCode": "AD",
	        "followers": "220"
	    },
	    --------------- many lines later -----------------
	    {
	        "CCode": "XK",
	        "followers": "2207"
	    },
	    {
	        "CCode": "ZA",
	        "followers": "21308"
	    }
	]

```

### /api/getdatapred/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:DumpTablePrediction:api\datalisting.go**

This function uses polynomial regression to produce a 3 coefficants that can be used to predict a dataset,

for example, if I wanted to generate a graph to predict the GDP to year:
`/api/getdatapred/gdp/year/GDP`

and I would get back

```json

	[8.110694932491328e+08,-837975.4864626579,216.51286834986442]

```


### /api/getcsvdata/:id/:x/:y
**GET - Auth Cookie Checked**

**Handled in api:GetCSV:api/datalisting.go**

This function was part of a old graph type that has now been converted.
At the time of writing it still works:

`/api/getcsvdata/gdp/year/GDP` (get a csv of year vs gdp value from the dataset gdp)

```

	"name","word","count"
	"1948","1948",276458
	"1949","1949",286752
	"1950","1950",297063
	"1951","1951",306281
	"1952","1952",307280

```

### /api/getreduceddata/:id
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

This call works like "getdata" but will return a limited set of it to ensure that the client requesting isnt overwelmed by a massive dataset.

### /api/getreduceddata/:id/:persent
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

This extends the functionallity of `/api/getreduceddata/:id` but offers you to select what persentage of the dataset you want.

### /api/getreduceddata/:id/:persent/:min
**GET - Auth Cookie Checked**

**Handled in api:DumpReducedTable:api/datalisting.go**

This extends the functionallity of `/api/getreduceddata/:id/:persent` but offers you to select what the min amount of data you want.

### /api/setdefaults/:id
**POST - Auth Cookie Checked**

**Handled in api:SetDefaults:api/**

Posting a string to this will put it next to that dataset.

### /api/getdefaults/:id
**GET - Auth Cookie Checked**

**Handled in api:GetDefaults:api/**

Gets the stings that was put into a dataset using `/api/setdefaults/:id`

### /api/identifydata/:id
**GET - Auth Cookie Checked**

**Handled in api:IdentifyTable:api/**

This parses the SQL table and outputs the parsed output in a friendly output for the end application to use to make desisions:

To get the layout of the `GDP` dataset then you can query `/api/identifydata/GDP`

```json

	{
	    "Cols": [
	        {
	            "Name": "year",
	            "Sqltype": "int"
	        },
	        {
	            "Name": "GDP",
	            "Sqltype": "float"
	        },
	        {
	            "Name": "change",
	            "Sqltype": "float"
	        },
	        {
	            "Name": "gdpindex",
	            "Sqltype": "float"
	        }
	    ],
	    "Request": "GDP"
	}

```

### /api/classifydata/:table/:col
**GET - Auth Cookie Checked**

**Handled in api:SuggestColType:api/**

This returns `true` or `false` for if the `:col` in dataset `:table` is numeric.

Example: `/api/classifydata/tweets/id`

```json
	
	true

```


### /api/stringmatch/:word
**GET - Auth Cookie Checked**

**Handled in api:FindStringMatches:api/**

This searches **all** datasets for a string inside, Useful for a related datasets bar.

`/api/stringmatch/BBC/` gives back:

```json

	[
	    {
	        "Count": 1,
	        "Match": "imp015c046f5b98b36486ed0c794fe8e1d711b2e097_6fa3560e60b224c9f68"
	    },
	    {
	        "Count": 1,
	        "Match": "imp10d50e756b7e31907d86aa2eda09ae724b1fed31_dded73f34760686b725"
	    },
	    {
	        "Count": 1,
	        "Match": "impe88601a6cee601e9fea744def250d16f264c963e_dded73f34760686b725"
	    },
	    {
	        "Count": 6,
	        "Match": "proc"
	    }
	]

```

### /api/stringmatch/:word/:x
**GET - Auth Cookie Checked**

**Handled in api:FindStringMatches:api/**

This function inherits the functioanlity of `/api/stringmatch/:word` but allows you to limit what `col` it came from using `:x`

### /api/relatedstrings/:guid
**GET - Auth Cookie Checked**

**Handled in api:GetRelatedDatasetByStrings:api/**

**Warning** This call is really slow because of the nature of it looking though very large amounts of data, use with care.

This looks at a dataset and will give you tables that contain similar data in them:

Example: `/api/relatedstrings/1a0408b99c6355ac28f137f7acd82363b4ea5524d3161ab9105b706b889`


```json

	[
	    {
	        "Match": "HEREFORD WORCESTER FIRE BRIGADE",
	        "Tables": [
	            "027b42a92100fed4347ff8cc41aeda13ea32be1ddded73f34760686b725",
	---many---
	            "ff66d78bea6e30cc3582473bcb08a37c38c4e5773759982115eec0d6048",
	            "fffdb25ee6b6a3f8e1fca12d4c6be865ec09b846c37b20778b8bac50d96"
	        ]
	    },
	    {
	        "Match": "HALTON BC",
	        "Tables": [
	            "027b42a92100fed4347ff8cc41aeda13ea32be1ddded73f34760686b725",
	            "02cc87162978680907a5c59ee76cfbde47357cf18c789d35dfcc847ddfd",
	---many---
	            "fffdb25ee6b6a3f8e1fca12d4c6be865ec09b846c37b20778b8bac50d96"
	        ]
	    }
	]

```