define ['jquery', 'app/PGPatternMatcher', 'app/PGOLMap', 'app/PGOverviewCharts'], ($, PGPatternMatcher, PGOLMap, PGOverviewCharts) ->
  'use strict'
  map = null
  charts = null
  data = {dataset: [], patterns: {}}

  resetCharts = (srcData) ->
    console.log srcData
    data = {dataset: [], patterns: {}}
    for item in srcData
      do (item) ->
        for key of item.tags
          do (key) ->
            if data.patterns[key]
              item[key] or= if data.patterns[key] is 'intNumber' or data.patterns[key] is 'floatNumber' then 0 else 'void'
            else
              data.patterns[key] = PGPatternMatcher.getPattern item[key]
    data.dataset = srcData
    charts = new PGOverviewCharts 'dummy', data, '#charts'

  updateMap = (data) ->
    console.log data

  $ () -> 
    map = new PGOLMap '#mapContainer'
    charts = new PGOverviewCharts 'dummy', data, '#charts'

    $(map).bind 'update', (evt, data) -> updateCharts data.elements
    $(map).bind 'search', (evt, data) -> resetCharts data.elements
    $(charts).bind 'update', (evt, data) -> updateMap data

