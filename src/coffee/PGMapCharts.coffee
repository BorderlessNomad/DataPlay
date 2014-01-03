define ['jquery', 'crossfilter', 'd3', 'dc', 'app/PGOverviewCharts'], 
($, crossfilter, d3, dc, PGOverviewCharts) ->
  class PGMapCharts extends PGOverviewCharts

    processData: ->
      @keys = (entry for entry of @data.patterns when (entry isnt 'Long' and entry isnt 'Lat'))
      @cfdata = crossfilter @data.dataset
      @dimensions.push(@cfdata.dimension (d) -> d[key]) for key in @keys
      if @dimensions.length > 1
        for i in [0..@dimensions.length-1]
          do (i) =>
            group = x: @keys[i], y: 'all', type: 'count', dimension: @dimensions[i], group: null 
            group.group = @dimensions[i].group().reduceCount() 
            @groups.push group
        #console.log entry.group.all() for entry in @groups
      