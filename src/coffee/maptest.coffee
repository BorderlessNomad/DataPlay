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
    #console.log data
    map.updateItems data.elements, false

  redefineDatasetKey = (data, srcKey, tgtKey) ->
    if srcKey isnt tgtKey
      for entry in data
        do (entry) ->         
            entry[tgtKey] = entry[srcKey]
            delete entry[srcKey]

  getDataSource = (guid) ->
    $.getJSON "/api/getdata/#{guid}", (data) ->
      if data.length
        patterns = {}
        for key of data[0]
          do (key) ->
            # Get the paterns for the key/value
            vp = PGPatternMatcher.getPattern data[0][key]
            kp = PGPatternMatcher.getKeyPattern key
            
            # Fix lat,lon keys for map
            switch kp
              when 'mapLongitude' 
                redefineDatasetKey data, key, 'lon'
                fixedKey = 'lon'
              when 'mapLatitude'
                redefineDatasetKey data, key, 'lat'
                fixedKey = 'lat'
              else
                fixedKey = key
             
            patterns[fixedKey] = valuePattern: vp, keyPattern: kp

            # Now parse ALL the data based on value pattern
            # TODO: Should lookup the key pattern before???
            entry[fixedKey] = PGPatternMatcher.parse(entry[fixedKey], vp) for entry in data

        chartWidth = $('#mapContainer').width()/4-2;
        chartHeight = $('#mapContainer').height()/6-2;
        # Generate Map charts and bind dc.js filtering events      
        charts = new PGMapCharts guid, {dataset: data, patterns: patterns}, '#charts', chartWidth, chartHeight
        # Event bindings to maps
        $(charts).bind 'update', (evt, data) -> updateMap data
        $(charts).bind 'updateOnlyItems', (evt, data) -> updateMapItems data

  $ () -> 
    # Generate Map -- should be after dataset load complete ????

    # TESTING: Leaflet map
    map = new PGLMap '#mapContainer'

    #map = new PGOLMap '#mapContainer'

    # Event bindings to charts
    $(map).bind 'update', (evt, data) -> updateCharts data if charts
    $(map).bind 'search', (evt, data) -> resetCharts data
    # Get data for the guid and create charts
    getDataSource guid

    

