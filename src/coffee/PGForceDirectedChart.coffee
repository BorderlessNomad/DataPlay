define ['d3', 'app/PGChart'], (d3, PGChart) ->
  class PGForceDirectedChart extends PGChart 
    forceDirectedEl: null
    force: null
    colors: null
    allNodes: null
    nodes: null
    nodeRadius: 20
    nodeClickTimeout: null
    links : null
    initialDepth: 2

    constructor: (container, margin, dataset, axes, patterns, limit, value, initialDepth) ->
      @value = value
      @initialDepth = initialDepth if initialDepth
      super container, margin, dataset, axes, patterns, limit
      

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->

    drawAxes: ->

    createForceDirected: ->
      @allNodes = []
      @flatten @allNodes, @dataset, 0
      @nodes = @filterNodes @allNodes 
      @forceDirectedEl = @chart.append('g')
        .attr('id', 'forcedirected')


    flatten: (nodes, node, depth) ->
      if node.values
        node.name = node.key
        node.children = node.values
        for child in node.values 
          do (child) =>
            child.parent = node
            @flatten nodes, child, depth+1
      node.depth = depth
      node.expand = depth <= @initialDepth
      nodes.push node

    filterNodes: (nodes) ->
      res = []
      for node in nodes 
        do (node) ->
          aux = node.parent
          aux = aux.parent while aux?.expand
          if not aux 
            node.children = if node.expand then node.values else null 
            res.push node  
      res

    initChart: ->
      #@drag = false
      super
      @colors = d3.scale.category10()
      @createForceDirected()
      @renderForceDirected()
      #d3.select('body').on('mouseup', () => @drag = true)

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->

    renderForceDirected: ->     
      that = @ 
      @links = d3.layout.tree().links(@nodes)

      @force = d3.layout.force()  
        .size([0.98*@width, 0.98*@height])  
        .nodes(@nodes)
        .links(@links)   
        .linkDistance(2*@nodeRadius)
        #.linkStrength(1)     
        #.friction(0)
        .charge(-300)
        #.alpha(0)
        .gravity(0.05)
        .start()

      links = @forceDirectedEl.selectAll('.node-link')
        .data(@links, (d) -> d.target.index)

      links.enter()
        .insert('line')
        .attr('class', 'node-link')

      links.transition()
        .style('stroke', '#999')
        .style('stroke-width', '3px')

      links.exit()
        .remove()

      drag = @force.drag()
        .on('dragstart', (d) -> d3.select(@).classed("fixed", d.fixed = true))
        #.on('drag', () => console.log d3.event.x)
        #.on('dragend', () => console.log 'dragend')

      nodes = @forceDirectedEl.selectAll(".node")
        .data(@nodes, (d) -> d.index)

      nodesEnter = nodes.enter()
        .append("g")
        .attr("class", "node")
        .on('mousedown', () => d3.event.stopPropagation())#@drag = false)
        .on('click', (d) =>
          clearTimeout @nodeClickTimeout
          if not d3.event.defaultPrevented
            @nodeClickTimeout = setTimeout(
              =>             
                # TODO: show options??
                d.expand = not d.expand
                @nodes = @filterNodes @allNodes
                @renderForceDirected()
              300
            )
        )
        .on('dblclick', (d) ->
          clearTimeout that.nodeClickTimeout
          d3.select(@).classed("fixed", d.fixed = false)
          d3.event.stopPropagation()
        )

      nodesEnter.append("circle")
        .attr('class', 'node-circle')
        .attr("r", @nodeRadius)
        .style("fill", (d) => @colors(d.depth))
        .style('stroke', '#000')

      nodesEnter.append('title')
        .attr('class', 'node-title')
        .text((d) => d.key ? d[@value])

      nodesEnter.append('image')
        .attr('class', 'node-img')
        .attr('xlink:href', '/img/icons/search-active.png')
        .attr('x', -0.6*@nodeRadius)
        .attr('y', -0.6*@nodeRadius)
        .attr('width', "#{1.2*@nodeRadius}px")
        .attr('height', "#{1.2*@nodeRadius}px")  

      nodes.call(drag)

      nodes.exit()
        .remove()

      @force.on('tick', (e) =>
        links.attr('x1', (d) -> d.source.x)
          .attr('y1', (d) -> d.source.y)
          .attr('x2', (d) -> d.target.x)
          .attr('y2', (d) -> d.target.y)
        nodes.attr('transform', (d) -> "translate(#{d.x}, #{d.y})")
      )

    updateChart: (dataset, axes) =>
      super dataset, axes
      @renderForceDirected
