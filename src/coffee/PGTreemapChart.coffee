define ['app/PGChart'], (PGChart) ->
  class PGTreemapChart extends PGChart 
    treemap: null
    value: null

    constructor: (container, margin, dataset, axes, patterns, limit, value) ->
      @value = value
      super container, margin, dataset, axes, patterns, limit
      

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->

    drawAxes: ->

    createTreemap: ->
      @treemap = @chart.append('g')
        .attr('id', 'treemap')

    initChart: ->
      super
      @createTreemap()
      @renderTreemap()

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->

    renderTreemap: ->
      colors = d3.scale.category20()

      treemap = d3.layout.treemap()
        .children((d) -> d.values)
        .value((d) => d[@value])
        #.value((d) => d.size)
        .round(true)
        .size([@width, @height])
        .sticky(false)
      
      #----------Crap tis to use the server dataset------------
      #@currDataset = json
      #-------------------------------------------------------
      
      nodes = treemap.nodes(@currDataset)
        .filter((d) -> not d.values)
      
      cells = @treemap.selectAll("g.cell")
        .data(nodes)

      cellEnter = cells.enter()
        .append("g")
        .attr("class", "cell")
      cellEnter.append("rect")
      cellEnter.append("text")

      cells.transition()
        .attr("transform", (d) -> "translate(#{d.x},#{d.y})")
        .select('rect')
        .attr('width', (d) -> d.dx-1)
        .attr('height', (d) -> d.dy-1)
        .style('fill', (d) -> colors(d.parent.key))

      cells.select('text')  
        .attr('x', (d) -> d.dx/2)
        .attr('y', (d) -> d.dy/2)
        .attr('dy', '.35em')
        .attr('text-anchor', 'middle')
        .text((d) => d[@value])
        .style('opacity', (d) ->
          d.w = @getComputedTextLength()
          if d.dx>d.w then 1 else 0
        )

      cells.exit()
        .remove()

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderTreemap()
