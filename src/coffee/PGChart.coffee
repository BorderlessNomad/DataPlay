 define ['jquery', 'd3', 'app/Common', 'app/PGChartPoint', 'app/PGChartPointView'], ($, d3, Common, PGChartPoint, PGChartPointView) -> 
  class PGChart 
    container: 'body'
    margin: {top: 50, right: 40, bottom: 50, left: 75}
    width: 800
    height: 600  
    dataset: []
    currDataset: []
    axes: {x: 'ValueX', y: 'ValueY'}
    axis: {x: null, y: null}
    patterns: null
    scale: {x: 1, y: 1}
    chart: null
    limit: 1000
    drag: true

    constructor: (container, margin, dataset, axes, patterns, limit) ->
      @container = container unless not container
      @margin = margin unless not margin
      @width = $(container).width() - @margin.left - @margin.right
      @height = $(container).height() - @margin.top - @margin.bottom    
      @dataset = dataset unless not dataset      
      @axes = axes unless not axes
      @patterns = patterns unless not patterns 
      @limit = limit unless not limit 
      @initChart()

    processDataset: -> 
      # TODO: Order first?????
      if @dataset.length > @limit
        inc = @dataset.length / @limit
        ind = 0
        @currDataset = []
        for i in [0..@dataset.length-1]
          do (i) =>
            if Math.floor(ind) is i
              @currDataset.push(@dataset[i])
              ind += inc
      else
        @currDataset = @dataset

    setScale: (axis, orient, ticks) ->
      pattern = @patterns[@axes[axis]].valuePattern
       # Get axes depending on axes pattern
      @scale[axis] = switch pattern
        when 'date' then d3.time.scale()
        when 'label' then d3.scale.ordinal()
        else d3.scale.linear()   
      switch pattern
        when 'label'
          dmn = []
          dmn.push d[0] for d in @currDataset when dmn.indexOf(d[0])<0
          rng = [0.01*@width, 0.98*@width]
          @scale[axis].domain(dmn)
            .rangeBands(rng)          
          @scale[axis].invert = (x) -> dmn[Math.round(dmn.length*(x-rng[0])/(rng[1]-rng[0]))]
        else 
          @scale[axis].domain(d3.extent(@currDataset, (d) -> d[if axis is 'x' then 0 else 1]))
            .range([
              if axis is 'x' then 0.01*@width else 0.98*@height
              if axis is 'x' then 0.98*@width else 0.01*@height
            ]) 
      @axis[axis] = d3.svg.axis()
        .scale(@scale[axis])
        .orient(orient ? 'bottom')
        .ticks(ticks ? 5)
      @axis[axis].tickFormat (d) =>
        switch pattern
          when 'date' then d.getFullYear()
          when 'percent' then "#{d}%"
          else d

    setScales: ->
      @setScale 'x', 'bottom', 5
      @setScale 'y', 'left', 8

    # getIvertedScale: (scale) ->
    #   switch @patterns[@axes.x]
    #     when 'label'
    #       dmn = scale.domain()
    #       rng = scale.range()
    #       (x) -> dmn[Math.round(dmn.length*(x-rng[0])/(rng[1]-rng[0]))]
    #     else scale.invert

    # --------------------- Chart creating Functions ------------------------ #
    drawChart: ->
      # The chart container
      svg = d3.select(@container).append("svg")
        .attr("width", @width  + @margin.left + @margin.right)
        .attr("height", @height + @margin.top + @margin.bottom)
        .on('mousedown', () -> d3.select(@).classed('panning', true))
        .on('mouseup', () -> d3.select(@).classed('panning', false))

      @chart = svg.append("g")
        .attr("transform", "translate(#{@margin.left},#{@margin.top})")

      # Clip chart
      @chart.append("clipPath")
         .attr('id', 'chart-area')
         .append('rect')
         .attr('x', 0)
         .attr('y', 0)
         .attr('width', @width)
         .attr('height', @height)

      svg.call d3.behavior.zoom().scaleExtent([1,10]).on 'zoom', () => 
          t = d3.event.translate
          s = d3.event.scale
          if @drag
            #svg.attr('-webkit-transform', "translate(#{t})scale(#{s})")
            svg.style('transform', "translate(#{t[0]}px,#{t[1]}px) scale(#{s})")
            svg.style('-webkit-transform', "translate(#{t[0]}px,#{t[1]}px) scale(#{s})")

    drawAxes: ->
      # Draw x axis
      @chart.append("g")
         .attr("class", "x axis")
         .attr("transform", "translate(0,#{@height})")
         .call(@axis.x)
         .append("text")
         .attr("id", "xLabel")
         .style("text-anchor", "end")
         .text(@axes.x)
         .attr("transform", "translate(#{@width},40)")
      # Draw y axis
      @chart.append("g")
        .attr("class", "y axis")
        .call(@axis.y)
        .append("text")
        .attr("id", "yLabel")
        .style("text-anchor", "middle")
        .text(@axes.y)  
        .attr("transform", "translate(0,-30)")

    initChart: ->
      # Preprocess dataset
      @processDataset()
      # Adjust scales to axis data types
      @setScales()
      # Draw the Chart
      @drawChart()
      # Draw the Axes
      @drawAxes()

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->
       # x axis transition
      @chart.select('.x.axis')
        .transition()
        .duration(1000)
        .call(@axis.x)
        .select("#xLabel")
        .text(@axes.x)
      # y axis transition
      @chart.select('.y.axis')
        .transition()
        .duration(1000)
        .call(@axis.y)
        .select("#yLabel")
        .text(@axes.y)

    updateChart: (dataset, axes) ->
      # Set new data values
      @dataset = dataset unless not dataset
      @axes = axes unless not axes
      # Preprocess dataset
      @processDataset()    
      # Adjust scales to axis data types
      @setScales()
      # Axes update /transitions
      @updateAxes()


    # ------------------ Events driven Functions --------------------- 
    newPointDialog: (x, y) ->
      # TODO: get id from server, saving it before, as it's new
      point = new PGChartPoint(id: 1)    
      #point.save()
      console.log point
      pointDialog = new PGChartPointView(
        model: point,
        # TODO: put here the real id
        id: "pointDialog#{point.get('id')}"
      )
      $(@container).append(pointDialog.render().$el)
