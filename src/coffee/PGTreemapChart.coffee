define ['d3', 'app/PGChart'], (d3, PGChart) ->
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
        .value((d) => Math.abs(d[@value]))
        #.value((d) => d.size)
        .round(true)
        .size([@width, @height])
        .sticky(false)
      
      nodes = treemap.nodes(@currDataset)
        .filter((d) -> d.depth is 1)
        #.filter((d) -> not d.values)
      
      #console.log nodes

      cells = @treemap.selectAll("g.cell")
        .data(nodes)

      cellsEnter = cells.enter()
        .append("g")
        .attr("class", "cell")
        .on('click', (d) => 
          if d.key
            #console.log d
            @currDataset = d
            @renderTreemap()
            # This will be the structure if needed request to server-side
            #@queryData [{key: ...., value: d.key}], (callback) -> ..... if d.key
          )
      cellsEnter.append("rect")
      cellsEnter.append('title')
      cellsEnter.append("text")

      cellsUpdate = cells.transition()
        .duration(1000)
        .attr("transform", (d) -> "translate(#{d.x},#{d.y})")
        .style('cursor', (d) -> if d.key then 'pointer' else 'default')
      cellsUpdate.select('rect')
        .attr('width', (d) -> d.dx-1)
        .attr('height', (d) -> Math.max(0,d.dy-1))
        .style('fill', (d, i) -> if d.value<0 then '#884444' else colors(i))
        #.style('fill', (d) -> colors(d.parent.key))
      cellsUpdate.select('title')
        .text((d) => if d.values then "#{d.key}\n#{d.value}" else "#{d.value}")
      cells.select('text')  
        .attr('x', (d) -> d.dx/2)
        .attr('y', (d) -> d.dy/2)
        .attr('dy', '.35em')
        .attr('text-anchor', 'middle')
        #.text((d) => d[@value])
        .text((d) => if d.key then d.key else d.value)
        #.style('opacity', (d) -> if d.dx>d.w then 1 else 0)
        .style('opacity', (d) -> 
          d.w = @getComputedTextLength()
          if d.dy>20 and d.dx>d.w then 1 else 0
        )
        
      cells.exit()
        .transition()
        .duration(1000)
        .remove()

    updateChart: (dataset, axes) =>
      super dataset, axes
      @renderTreemap()
