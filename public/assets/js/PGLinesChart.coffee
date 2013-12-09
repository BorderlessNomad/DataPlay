class window.PGLinesChart extends PGChart 

  # --------------------- Chart creating Functions ------------------------ #
  drawLines: ->
    that = @
    # Set chart line
    @line = d3.svg.line()
             #.interpolate("basis")
             .x((d) => @scale.x(d[0]))
             .y((d) => @scale.y(d[1]))
    # Draw chart line
    @chart.append('g')
      .attr('id', 'line01')
      .attr('clip-path', 'url(#chart-area)')
      .append("path")
      .attr("class", "line")
      .attr("d", @line(@currDataset))
      .on('mouseover', () ->
        $('.circle').remove()
        m = d3.mouse(@)
        console.log "#{that.scale.x.invert(m[0])},#{that.scale.y.invert(m[1])}"
        that.drawCircle(m[0], m[1])
      )
 
  initChart: ->
    super
    # Draw Chart lines
    @drawLines()

  # --------------------- Update Functions ------------------------ 
  updateLines: ->
    # Chart line transition
    @chart.selectAll(".line")
      .transition()
      .duration(1000)
      .attr("d", @line(@currDataset))

  updateChart: (dataset, axes) ->
    super dataset, axes
    # Curve update /transitions
    @updateLines()

  # ------------------ Events driven Functions --------------------- 
  drawCircle: (x, y) ->
    # Chart points as circles
    @chart.append("circle")        
      .attr("class", "circle")
      .attr("cx", x)
      .attr("cy", y)
      .attr("r", 5)
      .attr('fill', '#882244')
      .on('click', (d) => @newPointDialog(x, y))                  
      .append('title')
      .text((d) => "#{@axes.x}: #{@scale.x.invert(x)}\n#{@axes.y}: #{@scale.y.invert(y)}")


    # FIXME: Ask Ben for this Stuff!!!  
    # circles1.enter().attr('fill', "#00FF00");
    #window.GraphOBJ = circles1;
    #LightUpBookmarks();
    #fuckyoucoffescript.data().forEach(function(datum, i) { if (datum[1] > 500000 && datum[1] < 1000000) fuckyoucoffescript[0][i].setAttribute('fill','red'); })


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
