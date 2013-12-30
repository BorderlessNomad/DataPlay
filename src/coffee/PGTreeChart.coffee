define ['app/PGChart'], (PGChart) ->
  class PGTreeChart extends PGChart 
    treeEl: null
    tree: null
    diagonal: null
    i: 0
    value: null

    constructor: (container, margin, dataset, axes, patterns, limit, value) ->
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

    # --------------------- Update Functions ------------------------ 
    updateAxes: ->

    renderTree: ->
      @tree = d3.layout.tree()
        .children((d) -> d.values)
        .size([@height, @width])

      @diagonal = d3.svg.diagonal()
        .projection((d) -> [d.y, d.x])
      
      #----------Crap tis to use the server dataset------------
      #@currDataset = json
      #-------------------------------------------------------
      @currDataset.x0 = @height/2
      @currDataset.y0 = 0
      @render @currDataset
      
    render: (source) ->
      nodes = @tree.nodes(@currDataset).reverse()

      nodes.forEach((d) -> d.y = d.depth * 180)

      node = @treeEl.selectAll("g.node")
        .data(nodes, (d) => d.id or= ++@i)

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
        .style("fill", (d) -> if d._children then "lightsteelblue" else "#fff")

      nodeUpdate = node.transition()
        .attr("transform", (d) -> "translate(#{d.y},#{d.x})")

      nodeUpdate.select("circle")
        .attr("r", 4.5)
        .style("fill", (d) -> if d._children then "lightsteelblue" else "#fff")

      nodeExit = node.exit()
        .transition()
        .attr("transform", (d) -> "translate(#{source.y},#{source.x})")
        .remove()

      nodeExit.select("circle")
        .attr("r", 1e-6)

      nodeEnter.append("svg:text")
        .attr("x", (d) -> if d.values then -10 else 10)
        .attr("dy", ".35em")
        .attr("text-anchor", (d) -> if d.values then "end" else "start")
        .text((d) => if d.values then d.key else d[@value])
        #.text((d) => if d.values then d.key else d.size)
        .style("fill-opacity", 1e-6)

      nodeUpdate.select("text")
        .style("fill-opacity", 1)

      nodeExit.select("text")
        .style("fill-opacity", 1e-6)

      nodes.forEach((d) ->
          d.x0 = d.x
          d.y0 = d.y
      )

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

    toggle: (d) ->
      if d.values
        d._children = d.values
        d.values = null
      else
        d.values = d._children
        d._children = null

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderTree()
