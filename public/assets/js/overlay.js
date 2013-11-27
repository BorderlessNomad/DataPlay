'use strict';

function drawGraph(data, data2) {
	// Ok, xchart doesn't support double axes ...
	/*
	var width = $('#placeholder').width(),
		height = $('#placeholder').height();

	var x = d3.scale.linear().range([0, width]), 
		y = d3.scale.linear().range([height, 0]),
		color = d3.scale.category10(),
		xAxis = d3.svg.axis().scale(x).orient("bottom"),
		yAxis = d3.svg.axis().scale(y).orient("left"),
		yAxis2 = d3.svg.axis().scale(y).orient("right"),
		line = d3.svg.line().interpolate("basis")
				 .x(function(d) { return x(d[0]); })
    			 .y(function(d) { return y(d[1]); }),
   		svg = d3.select("#placeholder").append("svg")
    			.attr("width", width)
    			.attr("height", height)
  				.append("g")
    			.attr("transform", "translate(60)");

    	color.domain(d3.keys(data[0]);

    	var cities = color.domain().map(function(name) {
			return {
			  name: name,
			  values: data.map(function(d) {
			    return {date: d.date, temperature: +d[name]};
			  })
			};
		});
	

	Will follow tomorrow ......

	*/
	var options = {
		xScale: "linear",
		yScale: "linear",
		type: "line",
		main: [	
			{
		      	className: ".plot",
		      	type: "line",
		      	data: []
		  	},
			{
				className: ".plot2",
				type: "line",
		    	data: []
			}
		]
	};
	
	data.forEach(function(entry) {
		options.main[0].data.push({x: entry[0], y: entry[1]});
	});

	data2.forEach(function(entry) {
		options.main[1].data.push({x: entry[0], y: entry[1]});
	});

	//console.log(options);

	var myChart = new xChart('line', options, '#placeholder');

	var y2 = d3.scale.linear().range([$('#placeholder').height(), 0]),
		yAxisRight = d3.svg.axis().scale(y2).orient("right").ticks(5);  
	y2.domain([0, d3.max(data, function(d) { return Math.max(d[1]); })]); 
	d3.select('svg.xchart .scale').append("g")             
       .attr("class", "axis axisY")    
       .attr("transform", "translate(" + $('#placeholder').width()+100+" ,0)")   
       .style("fill", "red")       
       .call(yAxisRight);

}

function parseChartData(data, x, y) {
	var DataPool = [];
	for (var i = 0; i < data.length; i++) {
		var Unit = data[i];
		DataPool.push([ parseFloat(Unit[x]) , parseFloat(Unit[y]) ]);
	}
	return DataPool;
}

$( document ).ready(function() {
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

	$.getJSON( "/api/getinfo/" + guid, function( data ) {
		$('#FillInDataSet').html(data.Title);
		$(".wikidata").html(data.Notes);
	});

	function populatedKeys (Keys1,Keys2) {
		$("#pickx1axis").empty();
		$("#picky1axis").empty();
		$("#pickx2axis").empty();
		$("#picky2axis").empty();
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
		if (Keys1.length > 1) {
			$("#picky1axis").val(Keys1[1]);
		}
		if (Keys2.length > 1) {
			$("#picky2axis").val(Keys2[1]);
		}
	}

	var Keys2Populated = false;
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

		$.getJSON( "/api/getdata/" + localStorage['overlay1'], function( dataa ) {
			window.DataSet2 = dataa;

			if (dataa.length) {
				for(var key in dataa[0]) {
					Keys2.push(key);
				}
			}
			populatedKeys(Keys,Keys2);

			var DataPool2 = parseChartData(dataa, Keys2[0], Keys2[1]);
			drawGraph(DataPool,DataPool2);
			//$.plot("#placeholder", [ DataPool,DataPool2 ]);
		});
		
		
	});
	window.ReJigGraph = function() {
		var DataPool = parseChartData( window.DataSet, $("#pickx1axis").val(), $("#picky1axis").val());
		var DataPool2 = parseChartData( window.DataSet2, $("#pickx2axis").val(), $("#picky2axis").val()); 	

		drawGraph(DataPool,DataPool2);

		/*
		$.plot("#placeholder", [ 
			{ data: DataPool, label: $("#picky1axis").val() },
			{ data: DataPool2, label: $("#picky2axis").val(), yaxis: 2 }
		]);
		*/
	};
});
