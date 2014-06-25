define ['jquery', 'app/PGPatternMatcher'], ($, PGPatternMatcher) ->
  'use strict'
  class Common
    # TODO: Deprecated by PGPatternMatcher class
    @getPattern: (data) -> PGPatternMatcher.getPattern data

    @parseDate: (input) ->
      parts = input.split('-')
      new Date(parts[0], parts[1]-1, parts[2])

    @formatDate: (input) -> "#{input.getFullYear()}-#{input.getMonth()+1}-#{input.getDate()}"

    @parseAxisData: (data, axis) ->
      #console.log data, axis
      datapool = []
      # For the momet always force to get the pattern just in case data types change
      #pattern = axis.pattern ? PGPatternMatcher.getPattern data[0][axis.key]
      valuePattern = PGPatternMatcher.getPattern data[0][axis.key]
      for item in data
        do (item) =>
          datapool.push(PGPatternMatcher.parse item[axis.key], valuePattern)
      datapool

    @parseChartData: (data, x, y, guid, chartType, callback) ->
      if ['bars', 'pie', 'bubbles'].indexOf(chartType) >= 0
        $.get "/api/getdatagrouped/#{guid}/#{x.key}/#{y.key}", (dataset) =>
          @parseDataResults dataset, x, y, callback
      else
        @parseDataResults data, x, y, callback

    @parseDataResults: (data, x, y, callback) ->
      datapool = []
      xData = @parseAxisData(data, x)
      yData = @parseAxisData(data, y)
      datapool.push [xData[i] , yData[i]] for i in [0..data.length-1]
      callback datapool, {x: x.key, y: y.key}

    @quicksort: (dataset) ->
      if dataset.length <= 1
        dataset
      else
        pivot = Math.round(dataset.length/2)
        less = []
        greater = []
        for i in [0..dataset.length-1]
          do (i) ->
            if i isnt pivot
              if dataset[i][0]<=dataset[pivot][0] then less.push(dataset[i]) else greater.push(dataset[i])

        quicksort(less).concat([dataset[pivot]]).concat(quicksort greater)

    @saveUserDefaults: (guid, x, y) ->
      $.ajax(
        type: "POST"
        url: "/api/setdefaults/#{guid}"
        async: true
        data: "data=#{JSON.stringify({ x: x, y: y })}"
        success: (resp) -> #console.log(resp.result)
        error: (err) -> console.log(err)
      )

    @getUserDefaults: (guid, x, y, cb) ->
      $.getJSON "/api/getdefaults/#{guid}", (data) ->
        $(x).val data.x
        $(y).val data.y
        cb() if cb

    @getParameterByName: (name) ->
      name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]")
      regex = new RegExp("[\\?&]" + name + "=([^&#]*)")
      results = regex.exec(location.search)
      (if not results? then "" else decodeURIComponent(results[1].replace(/\+/g, " ")))
