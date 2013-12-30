define ['app/PGChart'], (PGChart) ->
  class PGEnclosureChart extends PGChart 
    enclosure: null
    r: 200
    value: null

    constructor: (container, margin, dataset, axes, patterns, limit, value) ->
      @value = value
      super container, margin, dataset, axes, patterns, limit
      

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->

    drawAxes: ->

    createEnclosure: ->
      @enclosure = @chart.append('g')
        .attr('id', 'enclosure')
        #.attr('transform', "translate(#{@r},#{@r})")

    initChart: ->
      super
      @r = Math.min(@width, @height)
      @createEnclosure()
      @renderEnclosure()

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->

    renderEnclosure: ->

      pack = d3.layout.pack()
        .size([@r, @r])
        .children((d) => d.values)
        .value((d) => Math.abs(d[@value]))
        #.value((d) => d.size)

      nodes = pack.nodes(@currDataset)
      #nodes = pack.nodes(json)

      console.log nodes

      circles = @enclosure.selectAll("circle")
        .data(nodes)

      circles.enter()
        .append("svg:circle")

      circles.transition()
        .attr("class", (d) -> if d.values then "parent" else "child")
        .attr("cx", (d) -> d.x)
        .attr("cy", (d) -> d.y)
        .attr("r", (d) -> d.r)
        .style('fill', (d) => if d[@value]<0 then '#884444' else '#448844')

      circles.exit()
        .transition()
        .attr("r", 0)
        .remove()

      labels = @enclosure.selectAll("text")
        .data(nodes)

      labels.enter()
        .append("svg:text")
        .attr("dy", ".35em")
        .attr("text-anchor", "middle")
        .style("opacity", 1);

      labels.transition()
        .attr("class",  (d) -> if d.values then "parent" else "child")
        .attr("x",  (d) -> d.x)
        .attr("y",  (d) -> d.y)
        .text((d) -> d.key)
        .style("opacity",  (d) -> if d.r>2 then 1 else 0.2)

      labels.exit()
        .remove()

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderEnclosure()
