class window.PGBarsChart extends PGChart 
  bars: null
  padding: 20

  # --------------------- Chart creating Functions ------------------------ #
  setScales: ->
    m = []
    m.push d[0] for d in @currDataset
    @scale.x = d3.scale.ordinal()
      .domain(m)
      .rangeBands([0.01*@width, 0.98*@width])

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

    @axis.x.tickFormat (d) -> d

    @axis.y = d3.svg.axis()
              .scale(@scale.y)
              .orient("left")

    @axis.y.tickFormat (d) =>
      switch DataCon.patterns[@axes.y]
        when 'date' then formatDate(d)
        when 'percent' then "#{d}%"
        else d

  createBars: ->
    @bars = @chart.append('g')
      .attr('id', 'bars')
      .attr('clip-path', 'url(#chart-area)')

  initChart: ->
    super
    @createBars()
    @renderBars()

  # --------------------- Update Functions ------------------------ 
  renderBars: ->
    bars = @bars.selectAll('rect.bar')
      .data(@currDataset)

    bars.enter()
      .append("rect")
      .attr("class", "bar")
      .on('click', (d) => @newPointDialog(d[0], d[1]))

    bars.exit()
      .transition()
      .duration(1000)
      .attr("x", -100)
      .remove()

    bars.transition()
      .duration(1000)
      .attr("x", (d) => @scale.x(d[0]))
      .attr("y", (d) => @scale.y(d[1]))
      .attr("width", (d) => Math.floor((@width/@currDataset.length)-@padding/@currDataset.length))
      .attr("height", (d) => @height-@scale.y(d[1]))

  updateChart: (dataset, axes) ->
    super dataset, axes
    @renderLines()

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
