define ['jquery', 'crossfilter', 'd3', 'dc'], ($, crossfilter, d3, dc) ->
  class PGOverviewCharts 
    guid: null
    container: 'body'
    data: null
    keys: []
    cfdata: null
    dimensions: []
    dimensionsMap: {}
    groups: []
    charts: [
      { id: 'rows', maxEntries: 15 }
      { id: 'bars', maxEntries: 30 }
      { id: 'pie', maxEntries: 45 }
      { id: 'bubbles', maxEntries: 60 }
      { id: 'line' }
    ]
    width: 238
    height: 120

    constructor: (@guid, @data, @container, width, height) ->
      @width = width if width
      @height = height if height
      @processData()
      @drawCharts()

    processData: ->
      @keys = (entry for entry of @data.patterns)
      @cfdata = crossfilter @data.dataset#?.slice(0, 10) # TODO: remove this limiting slice
      @dimensions.push(@cfdata.dimension (d) -> d[key]) for key in @keys
      if @dimensions.length > 1
        for i in [0..@dimensions.length-2]
          do (i) =>
            for j in [i+1..@dimensions.length-1]
              do (j) =>
                @addGroup i, j
                @addGroup j, i
        #console.log entry.group.all() for entry in @groups

    addGroup: (i, j) ->    
      xKey = @keys[i]
      xKeyPattern = @data.patterns[xKey].keyPattern
      xPattern = @data.patterns[xKey].valuePattern
      yKey = @keys[j]
      yKeyPattern = @data.patterns[yKey].keyPattern
      yPattern = @data.patterns[yKey].valuePattern

      #console.log "#{xKey}-#{xKeyPattern}-#{xPattern} \n #{yKey}-#{yKeyPattern}-#{yPattern}"
      # TESTING: whether to include a group or not into a chart
      # useGroup = switch xPattern
        # when 'excluded' then false
        # else true
      #   # TODO: insert more patterns here ...
      #   when 'label', 'date', 'postCode', 'creditCard' then true
      #   when 'intNumber' then switch xKeyPattern
      #     # TODO: insert key patterns here ... TODO
      #     when 'identifier', 'date' then true
      #     else false
      #   when 'floatNumber' then switch xKeyPattern
      #     when 'coefficient' then true
      #     else false
      #   else false
      # console.log useGroup
      useGroup = true # TODO: Crap this when using patterns above ...

      if useGroup
        # TODO: mark here if there's already a count, it has no sense to do more than once ...
        # useCount = .....
        group = x: xKey, y: yKey, type: 'count', dimension: @dimensions[i], group: null 
        group.group = @dimensions[i].group().reduceCount((d) -> d[yKey]) 
        @groups.push group

        useSum = switch yPattern
          # TODO: discard more patterns here ....
          when 'label', 'date', 'postCode', 'creditCard' then false
          when 'intNumber' then switch yKeyPattern
            when 'identifier', 'date' then false
            else true
          else true
        #console.log useSum

        if useSum
          group2 = x: xKey, y: yKey, type: 'sum', dimension: @dimensions[i], group: null
          group2.group = @dimensions[i].group().reduceSum((d) -> d[yKey]) 
          @groups.push group2

    drawLineChart: (entry, fixedId, xScale) ->
      chart = dc.lineChart "##{fixedId}"
      chart.width(@width)
        .height(@height)
        .margins({top: 10, right: 10, bottom: 30, left: 30})
        .dimension(entry.dimension)
        .group(entry.group)
        .transitionDuration(500)
        .elasticY(true)
        .x(xScale)           
        .xAxis()
        .ticks(3)
        .tickFormat((d) => if @data.patterns[entry.x] is 'date' then d.getFullYear() else d)
        # TODO: Everything should deliver a chart, thrash the workaround below when well-tested
        # if isNaN chart.yAxisMin()
        #   $(@container).find("##{entry.type}-#{entry.x}-#{entry.y}").remove()
        # else
        #   chart.yAxis().ticks(3)
      chart

    drawBarsChart: (entry, fixedId, xScale) ->
      chart = dc.barChart "##{fixedId}"
      chart.width(@width)
        .height(@height)
        .margins({top: 10, right: 10, bottom: 30, left: 30})
        .dimension(entry.dimension)
        .group(entry.group)
        .transitionDuration(500)
        .centerBar(true)  
        .gap(2)
        .elasticY(true)
        .x(xScale)           
        .xAxis()
        .ticks(3)
        .tickFormat((d) => if @data.patterns[entry.x] is 'date' then d.getFullYear() else d)
      chart

    drawRowsChart: (entry, fixedId) ->
      chart = dc.rowChart "##{fixedId}"
      chart.width(@width)
        .height(@height)
        .margins({top: 5, right: 10, bottom: 20, left: 10})
        .dimension(entry.dimension)
        .group(entry.group)
        .transitionDuration(500)
        .gap(1)
        .colors(d3.scale.category20())
        .label((d) -> d.key)
        .labelOffsetY(@height/(2*entry.group.size()))
        .title((d) -> d.value)
        .elasticX(true)
      chart

    drawPieChart: (entry, fixedId) ->
      chart = dc.pieChart "##{fixedId}"
      chart.width(@width)
        .height(@height)
        .radius(Math.min(@width, @height)/2)
        .innerRadius(0.1*Math.min(@width, @height))
        .dimension(entry.dimension)
        .group(entry.group)
        .transitionDuration(500)
        .colors(d3.scale.category20())
        .label((d) -> d.data.key)
        .minAngleForLabel(0.2)    
        .title((d) -> d.value)
      chart

    drawBubblesChart: (entry, fixedId, xScale) ->
      svg = d3.select("##{fixedId}")
        .append('svg')
        .attr('width', @width)
        .attr('height', @height)
      chart = dc.bubbleOverlay("##{fixedId}")
        .svg(svg)
        .width(@width)
        .height(@height)
        .dimension(entry.dimension)
        .group(entry.group)
        .transitionDuration(500)
        .keyAccessor((d) -> "Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_'))
        .colors(d3.scale.category20())
        .radiusValueAccessor((d) -> d.value)
        .r(d3.scale.linear().domain(d3.extent(entry.group.all(), (d) -> parseInt(d.value))))
        .maxBubbleRelativeSize(0.1)
        .minRadiusWithLabel(5) 
        .title((d) -> d.value)
      chart.point(
        "Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_')
        0.1*@width+0.8*xScale(d.key)
        0.2*@height+0.6*@height*Math.random()
      ) for d in entry.group.all()
      chart

    drawCharts: ->    
      lastCharts = []
      $(@container).html ''  
      for entry in @groups
        do (entry) =>
          switch @data.patterns[entry.x]
            # TODO: handle more patterns here .....
            when 'label'
              m = []
              m.push d.key for d in entry.group.all()
              #TODO@ this scale does not work well ... see weather UK ...
              xScale = d3.scale.ordinal()
                .domain(m)
                .rangeBands([0, @width])
            when 'date'
              xScale = d3.time.scale()
                #.domain(d3.extent(entry.group.all(), (d) -> PGPatternMatcher.parse(d.key, 'date')))
                .domain(d3.extent(entry.group.all(), (d) -> d.key))
                .range([0, @width])
            else
              xScale = d3.scale.linear()
                .domain(d3.extent(entry.group.all(), (d) -> parseInt(d.key)))
                .range([0, @width])
          fixedId = "#{entry.type}-#{entry.x}-#{entry.y}".replace(/[^a-zA-Z0-9_-]/gi, '_')
          $(@container).append """
            <div id='#{fixedId}'>
              <a href='/charts/#{@guid}?chart=lines&x=#{entry.x}&y=#{entry.y}'>
                <h4>#{entry.x}-#{entry.y}(#{entry.type})</h4>
              </a>
            <div>
          """
          chartId = null
          for dcChart in @charts
            do (dcChart) =>
              if not chartId and (not dcChart.maxEntries or entry.group.size()<dcChart.maxEntries) and lastCharts.indexOf(dcChart.id)<0
                chartId = dcChart.id
                lastCharts.push(dcChart.id)
          chart = switch chartId
            when 'rows' then @drawRowsChart entry, fixedId
            when 'bars' then @drawBarsChart entry, fixedId, xScale
            when 'pie' then @drawPieChart entry, fixedId
            when 'bubbles' then @drawBubblesChart entry, fixedId, xScale
            when 'line' then @drawLineChart entry, fixedId, xScale
            else @drawLineChart entry, fixedId, xScale

          # TESTING: How to get filtered data .... for maps or 3rd party elements
          chart.on "filtered", (chart, filter) =>
            console.log chart.dimension().top(Infinity)
            console.log filter
            # Trigger 'update' for focusing maps on items bounds and 'updateOnlyItems' for no focus
            $(@).trigger 'updateOnlyItems', {elements: chart.dimension().bottom Infinity}

          # TODO: CRAP this!! it causes all charts to trigger an event to the map
          # chart.on "postRedraw", (chart) =>
          #   #console.log chart.dimension().top(Infinity)
          #   $(@).trigger 'update', {elements: chart.dimension().bottom Infinity}

          lastCharts = [] if lastCharts.length is @charts.length 
          if ['bars', 'pie', 'bubbles'].indexOf(chartId)>-1
            urlChart = $("##{fixedId} a").attr('href').replace('lines',chartId) 
            $("##{fixedId} a").attr('href', urlChart) 
            
      dc.renderAll()
      resetAll = $("<div class='resetAll'><a class='btn btn-primary' role='button'>Reset All</a><div>")
      $(@container).prepend resetAll
      resetAll.click () =>
        dc.filterAll()
        dc.redrawAll()

      