express = require "express"
bodyParser = require "body-parser"
errorhandler = require "errorhandler"
jsonFile = require "json-file-plus"
path = require "path"
fs = require "fs"
hogan = require "hogan.js"
backend = path.join process.cwd(), 'backend.json'

app = express()

app.use bodyParser.urlencoded(extended: true)
app.use bodyParser.json()
app.use errorhandler()

port = process.env.PORT or 1937

router = express.Router()

router.route("/").get (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		res.json file.data

router.route("/").post (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to add." unless req.body?.ip?.length > 0

		index = file.data.indexOf req.body.ip

		return res.status(409).json error: "IP already exists!" if index isnt -1

		file.data.push req.body.ip

		file.save().then (->
			res.json file.data
		), (err) ->
			return res.status(500).json error: "Error while saving file."

router.route("/:ip").delete (req, res) ->
	jsonFile backend, (err, file) ->
		return res.status(500).json error: "Error while reading file." if err

		return res.status(400).json error: "No IP to remove." unless req.params?.ip?.length > 0

		index = file.data.indexOf req.params.ip

		return res.status(404).json error: "No such IP found!" if index is -1

		file.data.splice index, 1

		file.save().then (->
			res.json file.data
		), (err) ->
			return res.status(500).json error: "Error while saving file"

app.use '/', router

app.listen port, ->
	console.log "Express server listening on port %d in %s mode", port, app.settings.env
