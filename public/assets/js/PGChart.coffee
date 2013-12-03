_.templateSettings.variable = "data"

class window.PGChart 
  container: 'body'
  margin: {top: 50, right: 70, bottom: 50, left: 80}
  width: 800
  height: 600  
  dataset: []
  axes: {x: 'ValueX', y: 'ValueY'}
  axis: {x: null, y: null}
  scale: {x: 1, y: 1}
  chart: null

  constructor: (container, margin, dataset, axes) ->
    @container = container unless not container
    @margin = margin unless not margin
    @width = $(container).width() - @margin.left - @margin.right
    @height = $(container).height() - @margin.top - @margin.bottom    
    @dataset = dataset unless not dataset    
    @axes = axes unless not axes
    @initChart()

  initChart: () ->
    # Well, at the moment just one dataset

    @scale.x = d3.scale.linear()
          .domain(d3.extent(@dataset, (d) -> d[0]))
          .range([0, @width])
          .nice()
          #.clamp(true)

    @scale.y = d3.scale.linear()
          .domain([
            d3.min(@dataset, (d) -> d[1])
            d3.max(@dataset, (d) -> d[1])
          ])
          .range([@height, 0])
          .nice()

    @axis.x = d3.svg.axis()
              .scale(@scale.x)
              .orient("bottom")
              .ticks(5)

    @axis.y = d3.svg.axis()
              .scale(@scale.y)
              .orient("left")
      
    @line = d3.svg.line()
             #.interpolate("basis")
             .x((d) => @scale.x(d[0]))
             .y((d) => @scale.y(d[1]))

    @chart = d3.select(@container).append("svg")
            .attr("width", @width  + @margin.left + @margin.right)
            .attr("height", @height + @margin.top + @margin.bottom)
            .append("g")
            .attr("transform", "translate(#{@margin.left},#{@margin.top})")
    
    @chart.append("clipPath")
       .attr('id', 'chart-area')
       .append('rect')
       .attr('x', 0)
       .attr('y', 0)
       .attr('width', @width)
       .attr('height', @height)
    
    @chart.append("g")
       .attr("class", "x axis")
       .attr("transform", "translate(0,#{@height})")
       .call(@axis.x)
       .append("text")
       .attr("id", "xLabel")
       .style("text-anchor", "start")
       .text(@axes.x)
       .attr("transform", "translate(#{@width+20},0)")

    @chart.append("g")
      .attr("class", "y axis")
      .call(@axis.y)
      .attr("x", -20)
      .append("text")
      .attr("id", "yLabel")
      .attr("y", -30)
      .attr("dy", ".71em")
      .style("text-anchor", "end")
      .text(@axes.y)
    ###
    point = @chart.selectAll(".point")
               .data(@dataset)
               .enter().append("g")
               .attr("class", "point");
    ###
    @chart.selectAll("path")
       .data(@dataset)
       .enter().append("path")
       .attr("class", "line")
       .attr("d", (d) => @line(@dataset))
       .style("stroke", "#335577")

    # console.log(@dataset)

    circles1 = @chart.append('g')
      .attr('id', 'circles1')
      .attr('clip-path', 'url(#chart-area)')
      .selectAll(".circle1")
      .data(@dataset)
      .enter()
        .append("circle")        
        .attr("class", "circle1")
        .attr("cx", (d) => @scale.x(d[0]))
        .attr("cy", (d) => @scale.y(d[1]))
        .attr("r", 5)
        .attr('fill', '#882244')
        .on('mouseover', (d) -> 
          d3.select(@)
            .transition()
            .duration(500)
            .attr('r', 50)
            # .attr('fill', '#aa3377')
        )
        .on('mouseout', (d) ->
          d3.select(@)
            .transition()
            .duration(1000)
            .attr('r', 5)
            # .attr('fill', '#882244') 
        )
        .on('mousedown', (d) -> 

          SavePoint( d[0], d[1])
          d3.select(@)
            .transition()
            .duration(100)
            .attr('r', 5)
            .attr('fill', '#DEADBE')

          # @getPointData(1, @scale.x(d[0]), @scale.y(d[1]))
        )
    circles1
        .append('title')
        .text((d) => "#{@axes.x}: #{d[0]}\n#{@axes.y}: #{d[1]}")

    # circles1.enter().attr('fill', "#00FF00");
    window.GraphOBJ = circles1;
    LightUpBookmarks();
    #fuckyoucoffescript.data().forEach(function(datum, i) { if (datum[1] > 500000 && datum[1] < 1000000) fuckyoucoffescript[0][i].setAttribute('fill','red'); })
    ### This is for a second overlayed chart

    y2 = d3.scale.linear()
          .domain([
              d3.min(@dataset2, (d) -> d[1]),
              d3.max(@dataset2, (d) -> d[1])
            ])
          .range([height, 0])
          .nice()

    yAxis2 = d3.svg.axis()
               .scale(y2)
               .orient("right")

    line2 = d3.svg.line()
          #.interpolate("basis")
          .x((d) -> x(d[0]))
          .y((d) -> y2(d[1]))

    svg.append("g")
      .attr("class", "y axis")
      .attr("transform", "translate(" + @width + ",0)")
      .call(yAxis2)
      .append("text")
      .attr("y", 6)
      .attr("dy", ".71em")
      .style("text-anchor", "end")
      # TODO: Should be Key for Y Axis
      .text("Value2(m)")

    point2 = svg.selectAll(".point2")
                .data(@datasource2)
                .enter().append("g")
                .attr("class", "point2")
    point2.append('g')
        .attr('id', 'points2')
        .attr('clip-path', 'url(#chart-area)').append("path")
        .attr("class", "line")
        .attr("d", (d) -> line2(@dataset2))
        .style("stroke", "#775533")
      
    circles2 = svg.append('g')
                  .attr('id', 'circles2')
                  .attr('clip-path', 'url(#chart-area)')
                  .selectAll(".circle2")
                  .data(@dataset2)
                  .enter()
                  .append("circle")
                  .attr("class", "circle2")
                  .attr("cx", (d) -> x(d[0]))
                  .attr("cy", (d) -> y2(d[1]))
                  .attr("r", 5)
    ###


  updateChart: (dataset, axes) ->
    # Well, at the moment just one dataset

    @dataset = dataset unless not dataset    
    @axes = axes unless not axes

    @scale.x.domain(d3.extent(@dataset, (d) -> d[0]))

    @scale.y.domain([
      d3.min(@dataset, (d) -> d[1])
      d3.max(@dataset, (d) -> d[1])
    ])
   
    @chart.select('.x.axis')
      .transition()
      .duration(1000)
      .call(@axis.x)    
      .select("#xLabel")
      .text(@axes.x)  

    @chart.select('.y.axis')
      .transition()
      .duration(1000)
      .call(@axis.y)
      .select("#yLabel")
      .text(@axes.y)
        
    @chart.selectAll(".line")
      .data(@dataset)   
      .transition()
      .duration(1000) 
      .attr("d", (d) => @line(@dataset))
    
    @chart.selectAll(".circle1")
      .data(@dataset)   
      .transition()
      .duration(1000)    
      .attr("cx", (d) => @scale.x(d[0]))
      .attr("cy", (d) => @scale.y(d[1]))
      .select('title')
      .text((d) => "#{@axes.x}: #{d[0]}\n#{@axes.y}: #{d[1]}")

  getPointData: (id, x, y) ->
    point = new PGChartPoint(id: id)   
    jqxhr = point.fetch()
    jqxhr.success (model, response, options) =>
      console.log("Success!")
      console.log(model)
      console.log(response)
      console.log(options)
      @spawnPointDialog(model, x, y)
    jqxhr.error (model, response, options) =>
      console.log("Error!")
      console.log(model)
      console.log(response)
      console.log(options)
      @spawnPointDialog(point, x, y)


  spawnPointDialog: (point, x, y) ->
    pointTemplate = _.template $("#pointDataTemplate").html()
    pointInfoTemplate = pointTemplate point.toJSON()
    console.log pointInfoTemplate
    pointInfo = $(@container).first().append(pointInfoTemplate).find('.pointInfo').last()

    pointInfo.css('left', x).css('top', y)
    pointInfo.find('submit').click () ->
      title =  pointInfo.find('.pointInfoTitleInput').val()
      text =  pointInfo.find('.pointInfoTextInput').val()
      point.set {title: title, text: text}
      point.save()
      console.log(point)
      pointInfo.remove()


  # TODO: Neeeed for finish it!!!!!
    
