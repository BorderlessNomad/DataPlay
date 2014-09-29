express = require "express"
bodyParser = require "body-parser"
errorhandler = require "errorhandler"
jsonFile = require "json-file-plus"
path = require "path"
fs = require "fs"
swig = require "swig"
backend = path.join process.cwd(), 'backend.json'

app = express()

app.use bodyParser.urlencoded(extended: true)
app.use bodyParser.json()
app.use errorhandler()

port = process.env.PORT or 1937

compileTemplate = (data) ->
	for value, key in data.backends
		data.backends[key].id = "master#{key+1}"

	output = swig.renderFile "haproxy.cfg.template",
		backends: data.backends

	fs.writeFile "haproxy.cfg", output, (err) ->
		return err if err
		console.log "Done!"

router = express.Router()

router.route("/").get (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		res.json file.data.backends

router.route("/").post (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to add." unless req.body?.ip?.length > 0

		for value, key in file.data.backends
			if value.endpoint is req.body.ip
				return res.status(409).json error: "IP already exists!"

		file.data.backends.push
			endpoint: req.body.ip
			timestamp: Date.now()

		file.save().then (->
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
			console.log key
			if value.endpoint is req.params.ip
				index = key
				break

		return res.status(404).json error: "No such IP found!" if index is false

		console.log index

		file.data.backends.splice index, 1

		file.save().then (->
			compileTemplate file.data

			res.json file.data.backends
		), (err) ->
			return res.status(500).json error: "Error while saving file"

app.use '/', router

app.listen port, ->
	console.log "Express server listening on port %d in %s mode", port, app.settings.env
