define ['app/PGChart'], (PGChart) ->
  class PGBarsChart extends PGChart 
    bars: null
    padding: 30
    barsSet: []

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->
      patternX = @patterns[@axes.x]
      @barsSet.push(d[0]) for d in @currDataset when @barsSet.indexOf(d[0]) < 0
      console.log @barsSet
      @scale.x = d3.scale.ordinal()
        .domain(@barsSet)
        .rangeBands([0.01*@width, 0.98*@width])

      patternY = @patterns[@axes.y]#Common.getPattern @currDataset[0][1]
      @scale.y = switch patternY
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
        switch patternX
          when 'date' then d.getFullYear()
          when 'percent' then "#{d}%"
          else d

      @axis.y = d3.svg.axis()
                .scale(@scale.y)
                .orient("left")

      @axis.y.tickFormat (d) =>
        switch patternY
          when 'date' then d.getFullYear()
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
        .attr("width", (d) => Math.floor((@width/@barsSet.length)-@padding/@barsSet.length))
        .attr("height", (d) => @height-@scale.y(d[1]))

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderBars()
