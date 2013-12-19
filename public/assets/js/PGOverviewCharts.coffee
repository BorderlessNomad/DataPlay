class window.PGOverviewCharts 
  container: 'body'
  data: null
  cfdata: null
  dimensions: []
  groups: []
  charts: []

  constructor: (@data, @container) -> 
    @processData()
    @drawCharts()

  processData: ->
    @cfdata = crossfilter @data.dataset
    @dimensions.push(@cfdata.dimension (d) -> d[key]) for key in @data.keys
    for i in [0..@dimensions.length-2]
      do (i) =>
        for j in [i+1..@dimensions.length-1]
          do (j) =>
            @addGroup i, j
            @addGroup j, i
    console.log entry.group.all() for entry in @groups

  addGroup: (i, j) ->
    xKey = @data.keys[i]
    xPattern = @data.patterns[xKey]
    yKey = @data.keys[j]
    yPattern = @data.patterns[yKey]
    group = x: xKey, y: yKey, type: 'count', dimension: @dimensions[i], group: null 
    group.group = @dimensions[i].group().reduceCount((d) -> d[yKey]) 
    @groups.push group
    if yPattern isnt 'label'
        group2 = x: xKey, y: yKey, type: 'sum', dimension: @dimensions[i], group: null
        group2.group = @dimensions[i].group().reduceSum((d) -> d[yKey]) 
        @groups.push group2

  drawCharts: ->
    for entry in @groups
      do (entry) =>
        switch @data.patterns[entry.xKey]
          when 'label'
            m = []
            m.push d.key for d in entry.group.all()
            xScale = d3.scale.ordinal().domain(m)
          else
            xScale = d3.scale.linear().domain(d3.extent(entry.group.all(), (d) -> d.key))
              .range([0, 240])
        container = $(@container).append(
          "<div class='xs-col-3' id='#{entry.x}-#{entry.y}-#{entry.type}'><h4>#{entry.x}-#{entry.y}(#{entry.type})</h4><div>"
        )
        chart = dc.barChart "##{entry.x}-#{entry.y}-#{entry.type}"
        chart.width(240)
          .height(120)
          .margins({top: 10, right: 10, bottom: 30, left: 30})
          .dimension(entry.dimension)
          .group(entry.group)
          .transitionDuration(500)
          .centerBar(true)  
          .gap(2)
          #.filter([3, 5])
          .x(xScale)
          .elasticY(true)
          .xAxis()
          .ticks(3)
          .tickFormat((d) -> d)
        chart.yAxis()
          .ticks(3)


    dc.renderAll()
    