express = require "express"
bodyParser = require "body-parser"
errorhandler = require "errorhandler"
jsonFile = require "json-file-plus"
path = require "path"
fs = require "fs"
swig = require "swig"
sys = require "sys"
exec = require("child_process").exec
proxy = path.join process.cwd(), 'proxy.json'

app = express()

app.use bodyParser.urlencoded(extended: true)
app.use bodyParser.json()
app.use errorhandler()

port = process.env.PORT or 1937

puts = (error, stdout, stderr) -> sys.puts stdout

compileTemplate = (data) ->
	for value, key in data.gamification
		data.gamification[key].id = "gamification#{key+1}"

	for value, key in data.compute
		data.compute[key].id = "compute#{key+1}"

	timestamp = Date.now()
	output = swig.renderFile "haproxy.cfg.template",
		generatedOn: timestamp
		gamification: data.gamification
		compute: data.compute

	fs.writeFile "haproxy.cfg", output, (err) ->
		return err if err

		console.log "Successfully generated haproxy.cfg on #{new Date(timestamp)}"

		console.log "Copy /etc/haproxy/haproxy.cfg to /etc/haproxy/haproxy.cfg.#{timestamp}"
		exec "cp -rf /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.#{timestamp}", puts

		console.log "Replace old config file [Use some force if needed]"
		exec "cp -rf haproxy.cfg /etc/haproxy/haproxy.cfg", puts

		console.log "Reload HAProxy"
		exec "service haproxy reload", puts

router = express.Router()

router.route("/").get (req, res) ->
	jsonFile proxy, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		# compileTemplate file.data

		res.json file.data

router.route("/:type").get (req, res) ->
	jsonFile proxy, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No Type to remove." unless req.params?.type?.length > 0

		# compileTemplate file.data

		if file.data[req.params.type]? then res.json file.data[req.params.type] else res.status(400).json error: "Invalid Type specified."

router.route("/:type").post (req, res) ->
	jsonFile proxy, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No Type to specified." unless req.params?.type?.length > 0

		return res.status(400).json error: "Invalid Type specified." unless file.data[req.params.type]?

		console.log req.body
		return res.status(400).json error: "No IP to add." unless req.body?.ip?.length > 0

		for value, key in file.data[req.params.type]
			return res.status(409).json error: "IP already exists!" if value.endpoint is req.body.ip

		timestamp = Date.now()

		file.data[req.params.type].push
			endpoint: req.body.ip
			timestamp: timestamp

		file.save().then (->
			console.log "[#{new Date(timestamp)}] added new endpoint:", req.body.ip

			compileTemplate file.data

			res.json file.data[req.params.type]
		), (err) ->
			return res.status(500).json error: "Error while saving file."

router.route("/:type/:ip").delete (req, res) ->
	jsonFile proxy, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No Type to specified." unless req.params?.type?.length > 0

		return res.status(400).json error: "Invalid Type specified." unless file.data[req.params.type]?

		return res.status(400).json error: "No IP to remove." unless req.params?.ip?.length > 0

		for value, key in file.data[req.params.type]
			if value.endpoint is req.params.ip
				index = key
				break

		return res.status(404).json error: "No such IP found!" unless index?

		timestamp = Date.now()

		file.data[req.params.type].splice index, 1

		file.save().then (->
			console.log "[#{new Date(timestamp)}] removed endpoint:", req.params.ip

			compileTemplate file.data

			res.json file.data[req.params.type]
		), (err) ->
			return res.status(500).json error: "Error while saving file"

app.use '/', router

server = app.listen port, ->
	host = server.address().address
	port = server.address().port

	console.log "Express server listening on http://%s:%d in '%s' mode", host, port, app.settings.env
