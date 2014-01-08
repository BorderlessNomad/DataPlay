define ['jquery', 'crossfilter', 'd3', 'dc', 'app/PGOverviewCharts'], 
($, crossfilter, d3, dc, PGOverviewCharts) ->
  class PGMapCharts extends PGOverviewCharts

    processData: ->
      #@keys = (entry for entry of @data.patterns when (entry isnt 'lon' and entry isnt 'lat'))
      for entry of @data.patterns
        do (entry) =>
          vPattern = @data.patterns[entry].valuePattern
          kPattern = @data.patterns[entry].keyPattern
          addDim = switch vPattern
            when 'label', 'date', 'postCode', 'creditCard', 'currency' then true
            when 'intNumber' then switch kPattern
              when 'date' then true
              else false
            when 'floatNumber' then switch kPattern
              when 'coefficient', 'mapLatitude', 'mapLongitude' then true
              else false
            else false
          @keys.push entry if addDim
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
      @dimensionsMap['lon'].filter (d) -> data.left < d < data.right
      console.log @getFilteredDataset()
      dc.redrawAll()
      # Trigger 'update' for focusing maps on items bounds and 'updateOnlyItems' for no focus
      $(@).trigger 'updateOnlyItems', {elements: @getFilteredDataset()}
