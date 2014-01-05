importScripts '/requirejs/require.js'

requirejs.config({
  baseUrl: "/",
  paths: {
    jquery: "jquery/dist/jquery.min",
    underscore: "underscore/underscore-min",
    crossfilter: "crossfilter/crossfilter.min",
    kMeans: "lib/kMeans",
    app: "js"
  },
  shim: {
    'underscore': {
      exports: '_'
    },
    'crossfilter': {
      exports: 'crossfilter'
    }
  }
})

redefineDatasetKey = (data, srcKey, tgtKey) ->
  if srcKey isnt tgtKey
    for entry in data
      do (entry) ->         
          entry[tgtKey] = entry[srcKey]
          delete entry[srcKey]

self.addEventListener(
  'message'
  (e) ->
    data = e.data;
    self.postMessage 'msg': "Worker started on #{data.type}"  
    #kmeans = require 'kmeans'
    switch data.type

      when 'patterns'       
        require ['app/PGPatternMatcher'], (PGPatternMatcher) ->
          patterns = {}
          for key of data.data[0]
            do (key) ->
              # Get the paterns for the key/value
              vp = PGPatternMatcher.getPattern data.data[0][key]
              kp = PGPatternMatcher.getKeyPattern key
              
              # Fix lat,lon keys for map
              switch kp
                when 'mapLongitude' 
                  redefineDatasetKey data.data, key, 'lon'
                  fixedKey = 'lon'
                when 'mapLatitude'
                  redefineDatasetKey data.data, key, 'lat'
                  fixedKey = 'lat'
                else
                  fixedKey = key
               
              patterns[fixedKey] = valuePattern: vp, keyPattern: kp

              # Now parse ALL the data based on value pattern
              # TODO: Should lookup the key pattern before???
              entry[fixedKey] = PGPatternMatcher.parse(entry[fixedKey], vp) for entry in data.data

          msg = "Worker on patterns finished"
          res = 'dataset': data.data, 'patterns': patterns
          self.postMessage 'type': 'patterns', 'msg': msg, 'result': res
          self

      when 'kMeans'       
        require ['kMeans'], ->
          km = new kMeans K: 100      
          km.cluster data.vectors
          converged = false
          while km.step() and not converged
            do ->
              km.findClosestCentroids()
              km.moveCentroids()
              #console.log(km.centroids)
              converged = km.hasConverged()
          msg = "Worker on kMeans finished"
          res = 'centroids': km.centroids, 'clusters': km.clusters
          self.postMessage 'type': 'kMeans', 'msg': msg, 'result': res
          self

      else
        self.postMessage "Unknown work: #{data.type}"
        self
  false
)
