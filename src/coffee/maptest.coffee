define ['jquery', 'app/PGPatternMatcher', 'app/PGOLMap', 'app/PGOverviewCharts'], ($, PGPatternMatcher, PGOLMap, PGOverviewCharts) ->
  'use strict'
  map = null
  charts = null
  data = {dataset: [], patterns: {}}

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
              data.patterns[key] = keyPattern: PGPatternMatcher.getKeyPattern(key), valuePattern: PGPatternMatcher.getPattern(item[key])
    data.dataset = srcData
    charts = new PGOverviewCharts 'dummy', data, '#charts'

  updateMap = (data) ->
    console.log data
    item.lat = item.Lat for item in data.elements
    item.lon = item.Long for item in data.elements
    map.updateItems data.elements

  $ () -> 
    map = new PGOLMap '#mapContainer'
    #charts = new PGOverviewCharts 'dummy', data, '#charts'

    $(map).bind 'update', (evt, data) -> updateCharts data
    $(map).bind 'search', (evt, data) -> resetCharts data
    

    # TESTING: Particular test for map - UK Weather
    guid = 'weather_uk'
    $.getJSON "/api/getdata/#{guid}", (data) ->
      if data.length
        patterns = {}
        for key of data[0]
          do (key) ->
            vp = PGPatternMatcher.getPattern data[0][key]
            kp = PGPatternMatcher.getKeyPattern key
            patterns[key] = valuePattern: vp, keyPattern: kp
            # Now parse ALL the data
            # TODO get into account key pattern before parsing everything???
            entry[key] = PGPatternMatcher.parse(entry[key], patterns[key].valuePattern) for entry in data
            patterns[key].valuePattern = if key isnt 'Callsign' then 'excluded' else vp
        charts = new PGOverviewCharts guid, {dataset: data, patterns: patterns}, '#charts'
        $(charts).bind 'update', (evt, data) -> updateMap data

