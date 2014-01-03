define ['jquery', 'app/PGPatternMatcher', 'app/PGOLMap', 'app/PGMapCharts'],
($, PGPatternMatcher, PGOLMap, PGMapCharts) ->
  'use strict'
  map = null
  charts = null
  data = {dataset: [], patterns: {}}
  guid = window.location.href.split('/')[window.location.href.split('/').length - 1]

  # TODO: actually this only works with OSM overpass API searches
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
    charts = new PGMapCharts 'dummy', data, '#charts'

  updateCharts = (data) ->
    # TODO

  updateMap = (data) ->
    #console.log data
    map.updateItems data.elements

  redefineDatasetKey = (data, srcKey, tgtKey) ->
    for entry in data
      do (entry) ->
        if srcKey isnt tgtKey
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
            patterns[key] = valuePattern: vp, keyPattern: kp

            # Fix lat,lon keys for map
            switch kp
              when 'mapLongitude' then redefineDatasetKey data, key, 'lon'
              when 'mapLatitude' then redefineDatasetKey data, key, 'lat'
                
            # Now parse ALL the data based on value pattern
            # TODO: Should lookup the key pattern before???
            entry[key] = PGPatternMatcher.parse(entry[key], patterns[key].valuePattern) for entry in data

        # Generate Map charts and bind dc.js filtering events
        $('#charts').html ''
        charts = new PGMapCharts guid, {dataset: data, patterns: patterns}, '#charts'
        $(charts).bind 'update', (evt, data) -> updateMap data

  $ () -> 
    # Generate Map and bind search and update events 
    map = new PGOLMap '#mapContainer'
    $(map).bind 'update', (evt, data) -> updateCharts data
    $(map).bind 'search', (evt, data) -> resetCharts data
    # Get data for the guid and create charts
    getDataSource guid

    

