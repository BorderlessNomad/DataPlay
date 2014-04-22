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


### /api/user
**GET - Auth Cookie Checked**

**Handled in api:CheckAuth:api/datalisting.go**

### /api/visited
**GET - Auth Cookie Checked**

**Handled in api:GetLastVisited:api/tracking.go**

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
