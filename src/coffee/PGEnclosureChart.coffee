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
      colors = d3.scale.category10()
      pack = d3.layout.pack()
        .size([@r, @r])
        .children((d) => d.values)
        .value((d) => Math.abs(d[@value]))
        #.value((d) => d.size)

      nodes = pack.nodes(@currDataset)

      console.log nodes

      circles = @enclosure.selectAll("circle")
        .data(nodes)

      circles.enter()
        .append("svg:circle")
        .on('click', (d) => 
          console.log d
          @currDataset = d
          @renderEnclosure()
          #@queryData [{key: }]
        )
        .append('title')

      circles.exit()
        .transition()
        .duration(1000)
        .attr("r", 0)
        .remove()

      circles.transition()
        .duration(1000)
        .attr("class", (d) -> if d.values then "parent" else "child")
        .attr("cx", (d) -> d.x)
        .attr("cy", (d) -> d.y)
        .attr("r", (d) -> d.r)
        .style('fill', (d) => if d[@value]<0 then '#884444' else colors(d.depth))
        .select('title')
        .text((d) => if d.values then "#{d.key}\n#{d.value}" else "#{d.value}")
        

      labels = @enclosure.selectAll("text")
        .data(nodes)

      labels.enter()
        .append("svg:text")
        .attr("dy", ".35em")
        .attr("text-anchor", "middle")

      labels.exit()
        .transition()
        .duration(1000)
        .style("opacity", 0)
        .remove()

      labels.transition()
        .attr("class",  (d) -> if d.values then "parent" else "child")
        .attr("x",  (d) -> d.x)
        .attr("y",  (d) -> d.y)
        .text((d) -> (if d.values then d.key else d.value) if d.r>20)
        .style("opacity",  (d) -> if d.r>30 then 1 else 0.5)

      
    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderEnclosure()
