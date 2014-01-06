define ['jquery', 'app/PGPatternMatcher', 'app/PGOLMap', 'app/PGLMap','app/PGMapCharts'],
($, PGPatternMatcher, PGOLMap, PGLMap, PGMapCharts) ->
  'use strict'
  map = null
  charts = null
  data = {dataset: [], patterns: {}}
  guid = window.location.href.split('/')[window.location.href.split('/').length - 1]

  # TODO: actually this only works with OSM overpass API searches and OpenLayers
  resetCharts = (srcData) ->
    console.log srcData
    data = {dataset: [], patterns: {}}
    for item in srcData.elements
      do (item) ->
        for key of item.tags
          do (key) ->
            if data.patterns[key]
              item[key] or= switch data.patterns[key].valuePattern 
                when 'intNumber', 'floatNumber' then 0
                else 'void'
            else
              vp = PGPatternMatcher.getPattern item[key]
              kp = PGPatternMatcher.getKeyPattern key
              data.patterns[key] = keyPattern: kp, valuePattern: vp
    data.dataset = srcData
    $('#charts').html ''
    charts = new PGMapCharts 'dummy', data, '#charts'

  # updates charts data against map bounds
  updateCharts = (data) ->
    #console.log data
    # Conversion for leaflet bounds --> CRAP if using OpenLayers
    data = left: data.getWest(), right: data.getEast(), top: data.getNorth(), bottom: data.getSouth()

    charts.updateBounds data

  # updates map items and bounds
  updateMap = (data) ->
    #console.log data
    map.externalTrigger = true
    map.updateItems data.elements, true

  # updates only map items
  updateMapItems = (data) ->
    console.log data
    map.updateItems data.elements, false

  redefineDatasetKey = (data, srcKey, tgtKey) ->
    if srcKey isnt tgtKey
      for entry in data
        do (entry) ->         
            entry[tgtKey] = entry[srcKey]
            delete entry[srcKey]

  getDataSource = (guid) ->
    $.getJSON "/api/getdata/#{guid}", (dataset) ->
      if dataset.length
        patterns = {}

        worker = new Worker '/js/worker.js'
        worker.postMessage 'type': 'patterns', 'data': dataset
        worker.addEventListener(
          'message'
          (e) ->
            console.log e.data.msg
            switch e.data.type
              when 'kMeans'
                #console.log "Finished in: #{e.data.result.currentIteration} iterations"
                console.log e.data.result.centroids, e.data.result.clusters
                data.means = centroids: e.data.result.centroids, clusters: e.data.result.clusters
                generateCharts()
                worker.terminate()
              when 'patterns'
                data.dataset = e.data.result.dataset
                data.patterns = e.data.result.patterns

                # TESTING: kmeans over lat, lon
                fields = ['lat', 'lon']
                vectors = ([item[fields[0]] , item[fields[1]]] for item in data.dataset)
                worker.postMessage 'type': 'kMeans', 'vectors': vectors

                console.log data
          false
        )

  #TODO: it should be in the web worker ....
  insertClusterData = (data) ->
    ds = data.dataset
    dmc = data.means.clusters
    (ds[j].cluster = i for j in dmc[i]) for i in [0..dmc.length-1]
    data

  generateCharts = ->
    chartWidth = $('#mapContainer').width()/Object.keys(data.dataset[0]).length-2;
    chartHeight = $('#mapContainer').height()/8-2;
    # Process data clusters if any ....
    fixedData = if data.means then insertClusterData(data) else data

    # Generate Map charts and bind dc.js filtering events      
    charts = new PGMapCharts guid, fixedData, '#charts', chartWidth, chartHeight
    # Event bindings to maps
    $(charts).bind 'update', (evt, data) -> updateMap data
    $(charts).bind 'updateOnlyItems', (evt, data) -> updateMapItems data
    # Update chart items at startup
    updateMapItems {elements: charts.getFilteredDataset()}

  $ () -> 
    # Generate Map -- should be after dataset load complete ????

    # TESTING: Leaflet map
    map = new PGLMap '#mapContainer'

    # The OpenLayers one ...
    #map = new PGOLMap '#mapContainer'

    # Event bindings to charts
    $(map).bind 'update', (evt, data) -> updateCharts data if charts
    $(map).bind 'search', (evt, data) -> resetCharts data
    # Get data for the guid and create charts
    getDataSource guid
