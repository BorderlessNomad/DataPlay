class window.PGChart 
  container: 'body'
  margin: {top: 50, right: 40, bottom: 50, left: 75}
  width: 800
  height: 600  
  dataset: []
  currDataset: []
  axes: {x: 'ValueX', y: 'ValueY'}
  axis: {x: null, y: null}
  scale: {x: 1, y: 1}
  chart: null
  limit: 1000
  drag: true

  constructor: (container, margin, dataset, axes, limit) ->
    @container = container unless not container
    @margin = margin unless not margin
    @width = $(container).width() - @margin.left - @margin.right
    @height = $(container).height() - @margin.top - @margin.bottom    
    @dataset = dataset unless not dataset      
    @axes = axes unless not axes 
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
      console.log @currDataset
    else
      @currDataset = @dataset

  setScales: ->
     # Get axes depending on axes pattern
    @scale.x = switch DataCon.patterns[@axes.x]
      when 'date' then d3.time.scale()
      else d3.scale.linear()

    @scale.x.domain(d3.extent(@currDataset, (d) -> d[0]))
      .range([0.01*@width, 0.98*@width])

    @scale.y = switch DataCon.patterns[@axes.y]
      when 'date' then d3.time.scale()
      else d3.scale.linear()

    @scale.y.domain([
            d3.min(@currDataset, (d) -> d[1])
            d3.max(@currDataset, (d) -> d[1])
          ])
          .range([0.98*@height, 0.01*@height])

    @axis.x = d3.svg.axis()
              .scale(@scale.x)
              .orient("bottom")
              .ticks(5)

    @axis.x.tickFormat (d) =>
      switch DataCon.patterns[@axes.x]
        when 'date' then formatDate(d)
        when 'percent' then "#{d}%"
        else d

    @axis.y = d3.svg.axis()
              .scale(@scale.y)
              .orient("left")

    @axis.y.tickFormat (d) =>
      switch DataCon.patterns[@axes.y]
        when 'date' then formatDate(d)
        when 'percent' then "#{d}%"
        else d

  # --------------------- Chart creating Functions ------------------------ #
  drawChart: ->
    # The chart container
    svg = d3.select(@container).append("svg")
      .attr("width", @width  + @margin.left + @margin.right)
      .attr("height", @height + @margin.top + @margin.bottom)

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
          svg.attr('transform', "translate(#{t})scale(#{s})")

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
