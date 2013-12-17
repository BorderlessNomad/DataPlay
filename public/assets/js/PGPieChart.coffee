class window.PGPieChart extends PGChart 
  pie: null
  innerRadius: 0
  outerRadius: 200

  # --------------------- Chart creating Functions ------------------------ #
  createPie: ->
    @pie = @chart.append('g')
      .attr('id', 'pie')
      .attr('clip-path', 'url(#chart-area)')
      .attr('transform', "translate(#{@outerRadius},#{@outerRadius})")

  initChart: ->
    @processDataset()
    @drawChart()
    @createPie()
    @renderPie()

  # --------------------- Update Functions ------------------------ 
  renderPie: ->
    colors = d3.scale.category20()
    pie = d3.layout.pie()
      .sort((d) -> d[0])
      .value((d) -> d[1])

    arc = d3.svg.arc()
      .outerRadius(@outerRadius)
      .innerRadius(@innerRadius)

    slices = @pie.selectAll('path.arc')
      .data(pie(@currDataset))

    slices.enter()
      .append("path")
      .attr("class", "arc")
      .attr('fill', (d,i) -> colors(i))
      .on('click', (d) => @newPointDialog(d[0], d[1]))

    slices.exit()
      .transition()
      .duration(1000)
      .remove()

    slices.transition()
      .duration(1000)
      .attrTween "d", (d) -> 
        currArc = @currArc
        currArc or= startAngle: 0, endAngle: 0
        interpolate = d3.interpolate currArc, d
        @currArc = interpolate 1
        (t) -> arc interpolate t

    labels = @pie.selectAll('text.label')
      .data(pie(@currDataset))

    labels.enter()
      .append("text")
      .attr("class", "label")

    labels.exit()
      .transition()
      .duration(1000)
      .remove()

    labels.transition()
      .duration(1000)
      .attr('transform', (d) -> "translate(#{arc.centroid(d)})")
      .attr('dy', '.35em')
      .attr('text-anchor', 'middle')
      .text((d) -> d.data[0])

  updateChart: (dataset, axes) ->
    @processDataset() 
    @renderPie()
