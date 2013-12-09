class window.PGPanChart extends PGChart 
  maxEl: 100
  xStart: 0
  sortedDataset: []
  xPanRef: null

  constructor: (container, margin, dataset, axes, limit, maxEl) -> 
    @maxEl = maxEl unless not maxEl
    @maxEl = Math.min @maxEl, dataset.length
    super container, margin, dataset, axes, limit

  setCurrDataset: (xStart) ->
    @xStart = xStart unless not xStart
    @currDataset = @sortedDataset.slice(@xStart*@maxEl, (@xStart+1)*@maxEl)
    console.log @currDataset
  
  processDataset: ->
    # Sort dataset in x-axis
    @sortedDataset = quicksort(@dataset)
    # Initially focus on central data
    @setCurrDataset(Math.floor((@dataset.length/@maxEl)/2))

  # --------------------- Chart creating Functions ------------------------ #
  drawLines: ->
    # Set chart line
    @line = d3.svg.line()
             #.interpolate("basis")
             .x((d) => @scale.x(d[0]))
             .y((d) => @scale.y(d[1]))
    # Draw chart line
    @chart.append('g')
      .attr('id', 'lines1')
      .attr('clip-path', 'url(#chart-area)')
      .selectAll("path")
      .data(@currDataset)
      .enter().append("path")
      .attr("class", "line")
      .attr("d", (d) => @line(@currDataset))
      .style("stroke", "#335577")

  drawCircles: ->
    # Chart points as circles
    circles1 = @chart.append('g')
      .attr('id', 'circles1')
      .attr('clip-path', 'url(#chart-area)')
      .selectAll(".circle1")
      .data(@currDataset)       
      .enter()
      .append("circle")        
      .attr("class", "circle1")
      .attr("cx", (d) => @scale.x(d[0]))
      .attr("cy", (d) => @scale.y(d[1]))
      .attr("r", 5)
      .attr('fill', '#882244')
      .on('mouseover', (d) -> 
        d3.select(@)
          .transition()
          .duration(500)
          .attr('r', 50)
          .attr('fill', '#aa3377')
      )
      .on('mouseout', (d) ->
        d3.select(@)
          .transition()
          .duration(1000)
          .attr('r', 5)
          .attr('fill', '#882244') 
      )
      .on('click', (d) => @getPointData(1, @scale.x(d[0]), @scale.y(d[1])))                  
      .append('title')
      .text((d) => "#{@axes.x}: #{d[0]}\n#{@axes.y}: #{d[1]}")

  setChartEvents: ->
    $(@container).on 'mousedown', (e) => 
      console.log "drag started at #{e.clientX}"
      @xPanRef = e.clientX

    $(@container).on 'mousemove', (e) => 
      if @xPanRef
        console.log "dragged at #{e.clientX}"
        # Get new x axis reference
        xStart = Math.max 0, Math.floor(@maxEl*(e.clientX-@xPanRef)/@width)
        console.log "New start at #{xStart}"
        @setCurrDataset xStart
        # Adjust scales to axis data types
        @setScales()
        # Axes update /transitions
        @updateAxes()

    $(@container).on 'mouseup', (e) => 
      console.log "drag ended at #{e.clientX}"
      @xPanRef = null

  initChart: ->
    super
    # Draw Chart curve
    @drawLines()
    # Draw Chart points/circles
    #@drawCircles()
    # Set Chart Events
    @setChartEvents()

  # --------------------- Update Functions ------------------------ 
  updateLines: ->
    # Chart line transition
    @chart.selectAll(".line")
      .data(@currDataset)
      .transition()
      .duration(1000)
      .attr("d", (d) => @line(@currDataset))

  updateCircles: ->
    # Char points/circles transition
    @chart.selectAll(".circle1")
      .data(@currDataset)
      .transition()
      .duration(1000)
      .attr("cx", (d) => @scale.x(d[0]))
      .attr("cy", (d) => @scale.y(d[1]))
      .select('title')
      .text((d) => "#{@axes.x}: #{d[0]}\n#{@axes.y}: #{d[1]}")

  updateChart: (dataset, axes) ->
    super
    # Curve update /transitions
    @updateLines()
    # Points/Circles update /transitions
    @updateCircles()
