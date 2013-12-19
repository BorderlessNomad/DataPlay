'use strict'

window.DataCon or= {}

guid = window.location.href.split('/')[window.location.href.split('/').length - 1]
$ () ->
  $.getJSON "/api/getdata/#{guid}", (data) ->
    DataCon.dataset = data
    DataCon.keys = []
    if data.length
      DataCon.patterns = {}
      for key of data[0]
        do (key) ->
          DataCon.keys.push key 
          DataCon.patterns[key] = getPattern data[0][key]
      new PGOverviewCharts DataCon, '#charts'
