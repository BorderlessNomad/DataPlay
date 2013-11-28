'use strict';

window.DataCon || (window.DataCon = {});

function drawGraph(data, data2) {
	// Ok, xchart doesn't support double axes ...
	var margin = {top: 50, right: 70, bottom: 50, left: 50},
    	width = $('#placeholder').width() - margin.left - margin.right,
    	height = $('#placeholder').height() - margin.top - margin.bottom,
    	x = d3.scale.linear().range([0, width]),
 		y = d3.scale.linear().range([height, 0]),
 		y2 = d3.scale.linear().range([height, 0]),
 		color = d3.scale.category10(),
 		color2 = d3.scale.category10(),
 		xAxis = d3.svg.axis().scale(x).orient("bottom"),
 		yAxis = d3.svg.axis().scale(y).orient("left"),
 		yAxis2 = d3.svg.axis().scale(y2).orient("right"),
 		line = d3.svg.line().interpolate("basis")
    			 .x(function(d) { return x(d.x); })
    			 .y(function(d) { return y(d.y); }),
    	svg = d3.select("#placeholder").append("svg")
    			 .attr("width", width + margin.left + margin.right)
			     .attr("height", height + margin.top + margin.bottom)
			     .append("g")
			     .attr("transform", "translate(" + margin.left + "," + margin.top + ")");


  color.domain(d3.keys(data[0]));
  color2.domain(d3.keys(data2[0]));

  //console.log(data);

  console.log(color);
  
  var points = color.domain().map(function(name) {
	    return {
	      name: name,
	      values: data.map(function(d) {
	        return {x: d[0], y: d[1]};
	      })
	    };
	  }),
	  points2 = color2.domain().map(function(name) {
	    return {
	      name: name,
	      values: data2.map(function(d) {
	        return {x: d[0], y: d[1]};
	      })
	    };
	  });
  
  console.log(points);
  console.log(points2);

  x.domain(d3.extent(data, function(d) { return d[0]; }));

  y.domain([
    d3.min(points, function(c) { return d3.min(c.values, function(v) { return v.y; }); }),
    d3.max(points, function(c) { return d3.max(c.values, function(v) { return v.y; }); })
  ]);
  
  y2.domain([
    d3.min(points2, function(c) { return d3.min(c.values, function(v) { return v.y; }); }),
    d3.max(points2, function(c) { return d3.max(c.values, function(v) { return v.y; }); })
  ]);
  
  console.log(x);
  console.log(y);
  console.log(y2);

  svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," + height + ")")
      .call(xAxis);
	
  // TODO: get text from key .....   
  svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)
      .append("text")
      .attr("y", -25)
      .attr("dy", ".71em")
      .style("text-anchor", "end")
      .text("Value(m)");
  
  // TODO: get text from key .....    
  svg.append("g")
      .attr("class", "y axis")
      .attr("transform", "translate(" + width + ",0)")
      .call(yAxis2)
      .append("text")
      .attr("y", -25)
      .attr("dy", ".71em")
      .style("text-anchor", "start")
      .text("Value2(m)");

  var point = svg.selectAll(".point")
      .data(points)
      .enter().append("g")
      .attr("class", "point");

  var circles = svg.selectAll(".circle")
      .data(points)
      .enter().append("circle")
      .attr("cx", function(d) { return d.x })
      .attr("cy", function(d) { return d.y })
      .attr("r", 2)
      .attr("class", "circle");
      
  var point2 = svg.selectAll(".point2")
      .data(points2)
      .enter().append("g")
      .attr("class", "point2");

  var circles = svg.selectAll(".circle")
      .data(points)
      .enter().append("circle")
      .attr("cx", function(d) { return d.x })
      .attr("cy", function(d) { return d.y })
      .attr("r", 2)
      .attr("class", "circle");

  console.log(point);
  console.log(point2);
      
  point.append("path")
      .attr("class", "line")
      .attr("d", function(d) { return line(d.values); })
      .style("stroke", function(d) { return color(d.name); });
      
  point2.append("path")
      .attr("class", "line")
      .attr("d", function(d) { return line(d.values); })
      .style("stroke", function(d) { return color2(d.name); });
}

function updateGraph() {
	var DataPool = parseChartData( window.DataSet, $("#pickx1axis").val(), $("#picky1axis").val()),
		DataPool2 = parseChartData( window.DataSet2, $("#pickx2axis").val(), $("#picky2axis").val()); 	
	$('#placeholder').html('');
	drawGraph(DataPool,DataPool2);
}


$( document ).ready(function() {	
	$("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	$(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	DataCon.guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
	DataCon.guid2 = localStorage['overlay1'];
	var guid = DataCon.guid,
		guid2 = DataCon.guid2;

	$.getJSON( "/api/getinfo/" + guid, function( data ) {
		$('#FillInDataSet').html(data.Title);
		$(".wikidata").html(data.Notes);
	});

	function populatedKeys (Keys1,Keys2) {
		// Empty selections
		$("#pickx1axis").empty();
		$("#picky1axis").empty();
		$("#pickx2axis").empty();
		$("#picky2axis").empty();
		// Populate selections
		for (var i = 0; i < Keys1.length; i++) {
			$("#pickx1axis").append($("<option></option>").val(Keys1[i]).text(Keys1[i]));
		}
		for (var i = 0; i < Keys1.length; i++) {
			$("#picky1axis").append($("<option></option>").val(Keys1[i]).text(Keys1[i]));
		}
		for (var i = 0; i < Keys2.length; i++) {
			$("#pickx2axis").append($("<option></option>").val(Keys2[i]).text(Keys2[i]));
		}
		for (var i = 0; i < Keys2.length; i++) {
			$("#picky2axis").append($("<option></option>").val(Keys2[i]).text(Keys2[i]));
		}
		// Set different axes first
		if (Keys1.length > 1) {
			$("#picky1axis").val(Keys1[1]);
		}
		if (Keys2.length > 1) {
			$("#picky2axis").val(Keys2[1]);
		}
		// Get the last user preferences if any ...
		getUserDefaults(guid, $("#pickx1axis"), $("#picky1axis"));
		getUserDefaults(guid2, $("#pickx2axis"), $("#picky2axis"));
	}

	$.getJSON( "/api/getdata/" + guid, function( data ) {
		var Keys = [], Keys2 = [];
		//console.log(data);
		window.DataSet = data;

		if (data.length) {
			for(var key in data[0]) {
				Keys.push(key);
			}
		}

		var DataPool = parseChartData(data, Keys[0], Keys[1]);
		//console.log(DataPool);

		$.getJSON( "/api/getdata/" + guid2, function( dataa ) {
			window.DataSet2 = dataa;

			if (dataa.length) {
				for(var key in dataa[0]) {
					Keys2.push(key);
				}
			}
			populatedKeys(Keys,Keys2);

			var DataPool2 = parseChartData(dataa, Keys2[0], Keys2[1]);
			drawGraph(DataPool, DataPool2);
		});
		
		
	});

	window.ReJigGraph = function() {
		saveUserDefaults(guid, $("#pickx1axis").val(), $("#picky1axis").val());
		saveUserDefaults(guid2, $("#pickx2axis").val(), $("#picky2axis").val());
		updateGraph();
	};

	$(window).resize(function() {    
	    $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	    $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	    updateGraph();
	});
});
