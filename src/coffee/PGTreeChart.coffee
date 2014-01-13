define ['d3', 'app/PGChart'], (d3, PGChart) ->
  class PGTreeChart extends PGChart 
    treeEl: null
    tree: null
    diagonal: null
    i: 0
    value: null
    initialDepth: 1

    constructor: (container, margin, dataset, axes, patterns, limit, value, initialDepth) ->
      @value = value
      super container, margin, dataset, axes, patterns, limit
      

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->

    drawAxes: ->

    createTree: ->
      @treeEl = @chart.append('g')
        .attr('id', 'tree')

    initChart: ->
      super
      @createTree()
      @renderTree()

    expandNode: (node, depth) ->
      node.expand = true
      if --depth
        @expandNode(child, depth) for child in node.values? 
      else
        for child in node.values?
          do (child) ->
            child._children = child.values
            child.values = null

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->

    renderTree: ->
      @tree = d3.layout.tree()
        .children((d) -> d.values)
        .size([0.98*@height, 0.98*@width])

      @diagonal = d3.svg.diagonal()
        .projection((d) -> [d.y, d.x])

      @currDataset.x0 = 0.49*@height
      @currDataset.y0 = 0.01*@width
      @expandNode @currDataset, @initialDepth if @initialDepth
      @render @currDataset
      
    render: (source) ->
      nodes = @tree.nodes(@currDataset)
        .filter((d) -> 
          #console.log d
          aux = d.parent
          aux = aux.parent while aux?.expand
          if not aux 
            if not d.expand
              d._children or= d.values
              d.values = null
              console.log d.depth, d
            else
              d.forceEnter = not d.forceEnter
          not aux
        )
        .reverse()

      nodes.forEach (d) -> d.y = d.depth * 180

      node = @treeEl.selectAll("g.node")
        .data(nodes, (d) => d.id or= ++@i)

      #console.log nodes

      nodeEnter = node.enter()
        .append("svg:g")
        .attr("class", "node")
        .attr("transform", (d) -> "translate(#{source.y0},#{source.x0})")
        .on("click", (d) =>
            @toggle(d)           
            @render(d)
        )
      nodeEnter.append("svg:circle")
        .attr("r", 1e-6)
        .style("fill", (d) -> if d._children then "lightsteelblue" else "#ffffff")
      nodeEnter.append("svg:text")
        #.attr("x", (d) -> if d.values then -10 else 10)
        .attr("x", (d) -> if d.values or d._children then 10 else -10)
        .attr("dy", ".35em")
        #.attr("text-anchor", (d) -> if d.values then "end" else "start")
        .attr("text-anchor", (d) -> if d.depth then "start" else "end")
        #.text((d) => if d.values then d.key else d[@value])
        .text((d) => d.key ? d[@value])
        #.text((d) => if d.values then d.key else d.size)
        #.style("fill", '#ffffff')
        #.style("fill-opacity", 1e-6)

      nodeUpdate = node.transition()
        .attr("transform", (d) -> "translate(#{d.y},#{d.x})")
      nodeUpdate.select("circle")
        .attr("r", 4.5)
        .style("fill", (d) -> if d._children then "lightsteelblue" else "#ffffff")
      nodeUpdate.select("text")
        .style("fill-opacity", 1)

      nodeExit = node.exit()
        .transition()
        .attr("transform", (d) -> "translate(#{source.y},#{source.x})")
        .remove()
      nodeExit.select("circle")
        .attr("r", 1e-6)
      nodeExit.select("text")
        .style("fill-opacity", 1e-6)

      nodes.forEach (d) ->
          d.x0 = d.x
          d.y0 = d.y

      link = @treeEl.selectAll("path.link")
        .data(@tree.links(nodes), (d) -> d.target.id)

      link.enter()
        .insert("svg:path", "g")
        .attr("class", "link")
        .attr("d", (d) =>
            o = x: source.x0, y: source.y0
            @diagonal(source: o, target: o)
        )

      link.transition()
        .attr("d", @diagonal);

      link.exit()
        .transition()
        .attr("d", (d) =>
            o = x: source.x, y: source.y
            @diagonal(source: o, target: o)
        )
        .remove()

      # nodes.forEach (d) => 
      #   if d.depth is 1 and d.values
      #     console.log d
      #     d._children = d.values
      #     d.values = null
      #     @render d

    toggle: (d) ->
      if d.values
        d._children = d.values
        d.values = null
      else
        d.values = d._children
        d._children = null
      d.expand = not d.expand
      console.log d

    updateChart: (dataset, axes) =>
      super dataset, axes
      @renderTree()
