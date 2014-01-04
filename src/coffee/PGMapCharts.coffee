define ['jquery', 'crossfilter', 'd3', 'dc', 'app/PGOverviewCharts'], 
($, crossfilter, d3, dc, PGOverviewCharts) ->
  class PGMapCharts extends PGOverviewCharts

    processData: ->
      @keys = (entry for entry of @data.patterns when (entry isnt 'Long' and entry isnt 'Lat'))
      @cfdata = crossfilter @data.dataset
      for key in @keys
        do (key) =>
          dim = @cfdata.dimension (d) -> d[key]
          @dimensions.push dim
          @dimensionsMap[key] = dim
      if @dimensions.length > 1
        for i in [0..@dimensions.length-1]
          do (i) =>
            group = x: @keys[i], y: 'all', type: 'count', dimension: @dimensions[i], group: null 
            group.group = @dimensions[i].group().reduceCount() 
            @groups.push group
        #console.log entry.group.all() for entry in @groups

    updateBounds: (data) ->
      console.log data
      @dimensionsMap['lat'].filter (d) -> data.bottom < d < data.top
      console.log @dimensionsMap['lat'].bottom Infinity
      @dimensionsMap['lon'].filter (d) -> data.left < d < data.right
      console.log @dimensionsMap['lon'].bottom Infinity
      dc.redrawAll()
      # Trigger 'update' for focusing maps on items bounds and 'updateOnlyItems' for no focus
      $(@).trigger 'updateOnlyItems', {elements: @dimensionsMap['lon'].bottom Infinity}