define ['jquery', 'd3', 'app/PGLinesChart'], ($, d3, PGLinesChart) ->
  class PGAreasChart extends PGLinesChart
    areas: null

    createModeButtons: ->
      super
      $('button.mode').click () => @renderAreas()

    # --------------------- Chart creating Functions ------------------------ #
    createAreas: ->
      @areas = @chart.append('g')
        .attr("id", "areas")
        .attr('clip-path', 'url(#chart-area)')

    initChart: ->
      super
      @createAreas()
      @renderAreas()

    # --------------------- Update Functions ------------------------
    renderAreas: ->
      area = d3.svg.area()
        .interpolate(@mode)
        .x((d) => @scale.x(d[0]))
        .y0((d) => @scale.y(0))
        .y1((d) => @scale.y(d[1]))

      areas = @areas.selectAll('path.area')
        .data([@currDataset])

      areas.enter()
        .append('path')
        .attr('class', 'area')

      areas.exit()
        .remove()

      areas.transition()
        .duration(1000)
        .attr("d", (d) => area(d))

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderAreas()
