define ['jquery', 'd3', 'app/PGChart'], ($, d3, PGChart) ->

  class PGPieChart extends PGChart
    pie: null
    innerRadius: 20
    outerRadius: 200

    # --------------------- Chart creating Functions ------------------------ #
    setScales: ->

    drawAxes: ->

    createPie: ->
      @pie = @chart.append "g"
        .attr "id", "pie"
        .attr 'transform', "translate(#{@outerRadius},#{@outerRadius})"

    initChart: ->
      super
      @outerRadius = Math.min(@width, @height) / 2
      @createPie()
      @renderPie()

    # --------------------- Update Functions ------------------------
    updateAxes: ->

    renderPie: ->
      colors = d3.scale.category20()
      pie = d3.layout.pie()
        .sort (d) ->
          d[0]
        .value (d) ->
          d[1]

      arc = d3.svg.arc()
        .outerRadius @outerRadius
        .innerRadius @innerRadius

      # Arc slices
      arcs = @pie.selectAll "g.slice"
        .data pie
        .enter()
        .append "svg:g"
        .attr "class", "slice"

      arcs.append "svg:path"
        .attr "fill", (d, i) ->
          color i
        .attr "d", arc

      arcs.append "svg.text"
        .attr "transform", (d) =>
          d.outerRadius = @outerRadius + 50
          d.innerRadius = @innerRadius + 45
          "translate(" + arc.centroid(d) + ")"
        .attr "text-anchor", "middle"
        .style "fill", "Purple"
        .style "font", "bold 12px Arial"
        .text (d, i) ->
          dataSet[i].legendLabel

      arcs.filter (d) ->
          d.endAngle - d.startAngle > 0.2
        .append "svg:text"
        .attr "dy", ".35em"
        .attr "transform", (d) =>
          d.outerRadius = @outerRadius
          d.innerRadius = @innerRadius
          "translate(" + arc.centroid(d) + ")rotate(" + angle(d) + ")"
        .style "fill", "White"
        .style "font", "bold 12px Arial"
        .text (d) ->
          d.data.magnitude

      slices = @pie.selectAll "path.arc"
        .data pie @currDataset

      slices.enter()
        .append "path"
        .attr "class", "arc"
        .attr "fill", (d, i) ->
          colors i
        .on "click", (d) =>
          @newPointDialog d[0], d[1]

      slices.exit()
        .transition()
          .duration 1000
        .remove()

      slices.transition()
          .duration 1000
        .attrTween "d", (d) =>
          currArc = @currArc
          currArc or= startAngle: 0, endAngle: 0
          interpolate = d3.interpolate currArc, d
          @currArc = interpolate 1
          (t) ->
            arc interpolate t

      labels = @pie.selectAll "text.label"
        .data pie @currDatase

      labels.enter()
        .append "text"
        .attr "class", "label"

      labels.exit()
        .transition()
          .duration 1000
        .remove()

      labels.transition()
        .duration(1000)
        .attr "transform", (d) ->
          dAng = (d.startAngle + d.endAngle) * 90 / Math.PI
          lAng = dAng + if dAng > 180 then 90 else -90
          # if dAng > 2 && lAng > 2
          #   dAng = 2
          #   lAng = 2
          diffAng = (d.endAngle - d.startAngle) * 180 / Math.PI
          lScale = if diffAng > 1 then Math.min(diffAng / 9, 3) else 0
          "translate(#{arc.centroid(d)})rotate(#{lAng})scale(#{lScale})"
        .attr "dy", ".35em"
        .attr "text-anchor", "middle"
        .text (d) =>
            switch @patterns[@axes.x].valuePattern
              when 'date' then d.data[0].getFullYear()
              when 'label', 'text' then d.data[0].substring 0, 20
              else d.data[0]
      return

    updateChart: (dataset, axes) ->
      super dataset, axes
      @renderPie()
