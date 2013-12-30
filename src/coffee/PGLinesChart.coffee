define ['jquery', 'app/PGChart'], ($, PGChart) ->
  class PGLinesChart extends PGChart 
    lines: null
    mode: 'cardinal'
    modes: ['linear', 'linear-closed', 'step-before', 'step-after', 'basis', 'basis-open', 
      'basis-close', 'bundle', 'cardinal', 'cardinal-open', 'cardinal-closed', 'monotone']

    createModeButtons: ->
      $('#modes').html('')
      @modes.forEach (mode) ->
        $('#modes').append("<button class='mode btn btn-info' data='#{mode}'>#{mode}</button>");
      that = @    
      $('button.mode').click () -> 
        that.mode = $(@).attr('data')
        that.renderLines() 

    # --------------------- Chart creating Functions ------------------------ #
    createLines: ->
      @lines = @chart.append('g')
        .attr('id', 'lines')
        .attr('clip-path', 'url(#chart-area)')

    initChart: ->
      super
      @createModeButtons()
      @createLines()
      @renderLines()

    # --------------------- Update Functions ------------------------ 
    renderLines: ->
      that = @

      line = d3.svg.line()
        .interpolate(@mode)
        .x((d) => @scale.x(d[0]))
        .y((d) => @scale.y(d[1]))

      lines = @lines.selectAll('path.line')
        .data([@currDataset])

      lines.enter()
        .append("path")
        .attr("class", "line")
        .on('mouseover', () ->
          $('.circle').remove()
          m = d3.mouse(@)
          that.drawCircle(m[0], m[1])
        )

      lines.exit()
        .remove()

      lines.transition()
        .duration(1000)
        .attr("d", (d) -> line(d))

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderLines()

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
