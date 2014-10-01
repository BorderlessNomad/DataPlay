express = require "express"
bodyParser = require "body-parser"
errorhandler = require "errorhandler"
jsonFile = require "json-file-plus"
path = require "path"
fs = require "fs"
swig = require "swig"
sys = require "sys"
exec = require("child_process").exec
backend = path.join process.cwd(), 'backend.json'

app = express()

app.use bodyParser.urlencoded(extended: true)
app.use bodyParser.json()
app.use errorhandler()

port = process.env.PORT or 1937

puts = (error, stdout, stderr) -> sys.puts stdout

compileTemplate = (data) ->
	for value, key in data.backends
		data.backends[key].id = "master#{key+1}"

	timestamp = Date.now()
	output = swig.renderFile "haproxy.cfg.template",
		generatedOn: timestamp
		backends: data.backends

	fs.writeFile "haproxy.cfg", output, (err) ->
		return err if err

		console.log "Successfully generated haproxy.cfg on #{new Date(timestamp)}"

		console.log "Copy /etc/haproxy/haproxy.cfg to /etc/haproxy/haproxy.cfg.#{timestamp}"
		exec "cp -rf /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.#{timestamp}", puts

		console.log "Replace old config file"
		exec "cp -rf haproxy.cfg /etc/haproxy/haproxy.cfg", puts

		console.log "Reload HAProxy"
		exec "service haproxy reload", puts

router = express.Router()

router.route("/").get (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		# compileTemplate file.data

		res.json file.data.backends

router.route("/").post (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to add." unless req.body?.ip?.length > 0

		for value, key in file.data.backends
			return res.status(409).json error: "IP already exists!" if value.endpoint is req.body.ip

		timestamp = Date.now()

		file.data.backends.push
			endpoint: req.body.ip
			timestamp: timestamp

		file.save().then (->
			console.log "[#{new Date(timestamp)}] added new endpoint:", req.body.ip

			compileTemplate file.data

			res.json file.data.backends
		), (err) ->
			return res.status(500).json error: "Error while saving file."

router.route("/:ip").delete (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to remove." unless req.params?.ip?.length > 0

		index = false
		for value, key in file.data.backends
			if value.endpoint is req.params.ip
				index = key
				break

		return res.status(404).json error: "No such IP found!" if index is false

		timestamp = Date.now()

		file.data.backends.splice index, 1

		file.save().then (->
			console.log "[#{new Date(timestamp)}] removed endpoint:", req.params.ip

			compileTemplate file.data

			res.json file.data.backends
		), (err) ->
			return res.status(500).json error: "Error while saving file"

app.use '/', router

app.listen port, ->
	console.log "Express server listening on port %d in %s mode", port, app.settings.env
