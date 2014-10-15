// Generated by CoffeeScript 1.8.0
var app, backend, bodyParser, compileTemplate, errorhandler, exec, express, fs, jsonFile, path, port, puts, router, swig, sys;

express = require("express");

bodyParser = require("body-parser");

errorhandler = require("errorhandler");

jsonFile = require("json-file-plus");

path = require("path");

fs = require("fs");

swig = require("swig");

sys = require("sys");

exec = require("child_process").exec;

backend = path.join(process.cwd(), 'backend.json');

app = express();

app.use(bodyParser.urlencoded({
  extended: true
}));

app.use(bodyParser.json());

app.use(errorhandler());

port = process.env.PORT || 1937;

puts = function(error, stdout, stderr) {
  return sys.puts(stdout);
};

compileTemplate = function(data) {
  var key, output, timestamp, value, _i, _len, _ref;
  _ref = data.backends;
  for (key = _i = 0, _len = _ref.length; _i < _len; key = ++_i) {
    value = _ref[key];
    data.backends[key].id = "master" + (key + 1);
  }
  timestamp = Date.now();
  output = swig.renderFile("haproxy.cfg.template", {
    generatedOn: timestamp,
    backends: data.backends
  });
  return fs.writeFile("haproxy.cfg", output, function(err) {
    if (err) {
      return err;
    }
    console.log("Successfully generated haproxy.cfg on " + (new Date(timestamp)));
    console.log("Copy /etc/haproxy/haproxy.cfg to /etc/haproxy/haproxy.cfg." + timestamp);
    exec("cp -rf /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg." + timestamp, puts);
    console.log("Replace old config file");
    exec("cp -rf haproxy.cfg /etc/haproxy/haproxy.cfg", puts);
    console.log("Reload HAProxy");
    return exec("service haproxy reload", puts);
  });
};

router = express.Router();

router.route("/").get(function(req, res) {
  return jsonFile(backend, function(err, file) {
    if (err) {
      return res.status(500).json({
        error: "Error while reading file."
      });
    }
    return res.json(file.data.backends);
  });
});

router.route("/").post(function(req, res) {
  return jsonFile(backend, function(err, file) {
    var key, value, _i, _len, _ref, _ref1, _ref2;
    if (err) {
      return res.status(500).json({
        error: "Error while reading file."
      });
    }
    if (!(((_ref = req.body) != null ? (_ref1 = _ref.ip) != null ? _ref1.length : void 0 : void 0) > 0)) {
      return res.status(400).json({
        error: "No IP to add."
      });
    }
    _ref2 = file.data.backends;
    for (key = _i = 0, _len = _ref2.length; _i < _len; key = ++_i) {
      value = _ref2[key];
      if (value.endpoint === req.body.ip) {
        return res.status(409).json({
          error: "IP already exists!"
        });
      }
    }
    file.data.backends.push({
      endpoint: req.body.ip,
      timestamp: Date.now()
    });
    return file.save().then((function() {
      compileTemplate(file.data);
      return res.json(file.data.backends);
    }), function(err) {
      return res.status(500).json({
        error: "Error while saving file."
      });
    });
  });
});

router.route("/:ip")["delete"](function(req, res) {
  return jsonFile(backend, function(err, file) {
    var index, key, value, _i, _len, _ref, _ref1, _ref2;
    if (err) {
      return res.status(500).json({
        error: "Error while reading file."
      });
    }
    if (!(((_ref = req.params) != null ? (_ref1 = _ref.ip) != null ? _ref1.length : void 0 : void 0) > 0)) {
      return res.status(400).json({
        error: "No IP to remove."
      });
    }
    index = false;
    _ref2 = file.data.backends;
    for (key = _i = 0, _len = _ref2.length; _i < _len; key = ++_i) {
      value = _ref2[key];
      console.log(key);
      if (value.endpoint === req.params.ip) {
        index = key;
        break;
      }
    }
    if (index === false) {
      return res.status(404).json({
        error: "No such IP found!"
      });
    }
    console.log(index);
    file.data.backends.splice(index, 1);
    return file.save().then((function() {
      compileTemplate(file.data);
      return res.json(file.data.backends);
    }), function(err) {
      return res.status(500).json({
        error: "Error while saving file"
      });
    });
  });
});

app.use('/', router);

app.listen(port, function() {
  return console.log("Express server listening on port %d in %s mode", port, app.settings.env);
});