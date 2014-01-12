define ['jquery', 'app/PGChart'], ($, PGChart) ->
  class PGLinesChart extends PGChart 
    lines: null
    brush: null
    undo:null
    undoData: []
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

    createBrush: ->
      @brush = d3.svg.brush()
        .x(@scale.x)
        .y(@scale.y)
        .on("brushend", =>
          if not @brush.empty()
            if not @undo 
              @undo = @chart.append('image')
                .attr('class', 'undo')
                .attr('xlink:href', '/img/icons/undo.png')
                .attr('x', 30)
                .attr('y', -30)
                .attr('width', '48px')
                .attr('height', '48px')
                .on('click', =>
                  data = @undoData.pop()
                  if not @undoData.length
                    @chart.selectAll('.undo').remove()
                    @undo = null
                  @scale.x.domain data.x
                  @scale.y.domain data.y
                  @updateScope()
                )
            @undoData.push x:@scale.x.domain(), y:@scale.y.domain()  
            ext = @brush.extent()  
            @scale.x.domain [ext[0][0],ext[1][0]]
            @scale.y.domain [ext[0][1],ext[1][1]]
            @updateScope()
        )

      @chart.append("g")
        .attr("class", "brush")
        .call(@brush)
        .selectAll("rect")
        .attr("y", -@margin.top)
        .attr("height", @height+@margin.top)

    initChart: ->
      @drag = false
      super
      @createModeButtons()
      @createBrush()
      @createLines()
      @renderLines()
      

    # --------------------- Update Functions ------------------------ 
    updateBrush: ->
      @brush.clear().x(@scale.x).y(@scale.y)
      @chart.selectAll(".brush").call(@brush)

    updateScope: ->
      @updateBrush()      
      @updateAxes()
      @renderLines()

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
          m = d3.mouse(@)
          that.drawCircle(m[0], m[1])
        )
        .on('mouseout', () ->
          d3.select('.circle').remove()
        )

      lines.exit()
        .remove()

      lines.transition()
        .duration(1000)
        .attr("d", (d) -> line(d))

    updateChart: (dataset, axes) =>
      super dataset, axes
      @updateBrush()
      @renderLines()
      

    # ------------------ Events driven Functions --------------------- 
    drawCircle: (x, y) ->
      # Chart points as circles
      @chart.append("circle")        
        .attr("class", "circle")
        .attr("cx", x)
        .attr("cy", y)
        .attr("r", 10)
        .attr('fill', '#882244')
        .on('mousedown', (d) => @newPointDialog(x, y))                  
        .append('title')
        .text((d) => "#{@axes.x}: #{@scale.x.invert(x)}\n#{@axes.y}: #{@scale.y.invert(y)}")
