express = require "express"
bodyParser = require "body-parser"
errorhandler = require "errorhandler"
jsonFile = require "json-file-plus"
path = require "path"
fs = require "fs"
swig = require "swig"
sys = require "sys"
exec = require("child_process").exec
cluster = path.join process.cwd(), "cluster.json"
datakey = "cluster"

require("console-stamp") console, "[yyyy-mm-dd HH:MM:ss.l o]"

app = express()

app.use bodyParser.urlencoded(extended: true)
app.use bodyParser.json()
app.use errorhandler()

port = process.env.PORT or 1937

isFirstRun = true

puts = (error, stdout, stderr) -> sys.puts stdout

compileTemplate = (data) ->
	for value, key in data.cluster
		data.cluster[key].id = key

	if data.cluster.length > 1
		isFirstRun = false

	console.log "[Compile] - Prepare -", data

	timestamp = Date.now()
	output = swig.renderFile "pgpool.conf.template",
		generatedOn: timestamp
		cluster: data.cluster

	fs.writeFile "pgpool.conf", output, (err) ->
		console.log "[API] - GET - Error", err if err
		return err if err

		console.log "[Compile] - Generated -", "pgpool.conf on #{new Date(timestamp)}"

		exec "cp -rf /etc/pgpool-II-94/pgpool.conf /etc/pgpool-II-94/pgpool.conf.#{timestamp}", puts
		console.log "[Copy] - Backup -", "/etc/pgpool-II-94/pgpool.conf => /etc/pgpool-II-94/pgpool.conf.#{timestamp}"

		exec "cp -rf pgpool.conf /etc/pgpool-II-94/pgpool.conf", puts
		console.log "[Copy] - Replace -", "pgpool.conf => /etc/pgpool-II-94/pgpool.conf"

		if isFirstRun
			console.log "[Service] - Restart -", "Start"
			exec "systemctl restart pgpool-II-94", puts
			console.log "[Service] - Restart -", "Success"
		else
			console.log "[Service] - Reload -", "Start"
			exec "/usr/pgpool-9.4/bin/pgpool reload", puts
			console.log "[Service] - Reload -", "Success"

router = express.Router()

router.route("/").get (req, res) ->
	jsonFile cluster, (err, file) ->
		console.log "[API] - GET - Error", err if err
		return res.status(500).json error: "Error while reading file." if err

		# compileTemplate file.data
		console.log "[API] - GET -", req.headers['x-forwarded-for'] || req.connection.remoteAddress, "-", req.body.ip

		res.json file.data

router.route("/").post (req, res) ->
	jsonFile cluster, (err, file) ->
		console.log "[API] - POST - Error", err if err
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to add." unless req.body?.ip?.length > 0

		console.log "[API] - POST -", req.headers['x-forwarded-for'] || req.connection.remoteAddress, "-", req.body.ip

		for value, key in file.data[datakey]
			return res.status(409).json error: "IP already exists!" if value.endpoint is req.body.ip

		timestamp = Date.now()

		file.data[datakey].push
			endpoint: req.body.ip
			timestamp: timestamp

		file.save().then (->
			console.log "[API] - POST - Success", req.body.ip

			compileTemplate file.data

			res.json file.data[datakey]
		), (err) ->
			console.log "[API] - POST - Error", err
			return res.status(500).json error: "Error while saving file."

router.route("/:ip").delete (req, res) ->
	jsonFile cluster, (err, file) ->
		console.log "[API] - DELETE - Error", err if err
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to remove." unless req.params?.ip?.length > 0

		for value, key in file.data[datakey]
			if value.endpoint is req.params.ip
				index = key
				break

		console.log "[API] - DELETE -", req.headers['x-forwarded-for'] || req.connection.remoteAddress, "-", req.params.ip

		return res.status(404).json error: "No such IP found!" unless index?

		return res.status(400).json error: "Can not remove last node from Cluster, must have atleast one running node." if Object.keys(file.data[datakey]).length <= 1

		timestamp = Date.now()

		file.data[datakey].splice index, 1

		file.save().then (->
			console.log "[API] - DELETE - Success", datakey, req.params.ip

			compileTemplate file.data

			res.json file.data[datakey]
		), (err) ->
			console.log "[API] - DELETE - Error", err
			return res.status(500).json error: "Error while saving file"

app.use '/', router

server = app.listen port, ->
	host = server.address().address
	port = server.address().port

	console.log "[API] - INIT - [", app.settings.env, "]", "http://", host, ":", port
