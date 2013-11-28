'use strict';

function drawGraph(data) {
	var options = {
		xScale: "linear",
		yScale: "linear",
		type: "line-dotted",
		main: [	
			{
				className: ".plot",
		    	data: []
		  	}
		]
	};
	
	data.forEach(function(entry) {
		options.main[0].data.push({x: entry[0], y: entry[1]});
	});

	//console.log(options);
	var myChart = new xChart('line', options, '#placeholder');
}

function updateGraph() {
	var DataPool = parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()); 	
	$('#placeholder').html('');
	drawGraph(DataPool);
}

$( document ).ready(function() {
	$("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	$(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

	$.getJSON( "/api/getinfo/" + guid, function( data ) {
		$('#FillInDataSet').html(data.Title);
		$(".wikidata").html(data.Notes);
	});

	function populatedKeys (Keys) {
		$("#pickxaxis").empty();
		$("#pickyaxis").empty();
		for (var i = 0; i < Keys.length; i++) {
			$("#pickxaxis").append($("<option></option>").val(Keys[i]).text(Keys[i]));
		}
		for (var i = 0; i < Keys.length; i++) {
			$("#pickyaxis").append($("<option></option>").val(Keys[i]).text(Keys[i]));
		}
		if (Keys.length > 1) {
			$("#pickyaxis").val(Keys[1]);
		}	
		// Get the last user preferences if any ...
		getUserDefaults(guid, $("#pickxaxis"), $("#pickyaxis"));
	}

	$.getJSON( "/api/getdata/" + guid, function( data ) {
		//console.log(data);
		// so data is an array of shit.
		window.DataSet = data;
		var Keys = [];
		if (data.length) {
			for(var key in data[0]) {
				Keys.push(key);
			}
		}
		populatedKeys(Keys);
		drawGraph(parseChartData(data, Keys[0], Keys[1]));
	});

	window.ReJigGraph = function() {
		drawGraph(parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()));
	};
	$('#SetupOverlay').on("click",function() {
		// Okay so we need to put into local storage the current GUID
		// so when the overlay panal is called we can use it later on the
		// overlay page.
		localStorage['overlay1'] = guid;
		window.location.href = '/search/overlay';
	});

	window.ReJigGraph = function() {
		saveUserDefaults(guid, $("#pickxaxis").val(), $("#pickyaxis").val());
		updateGraph();
	};

	$(window).resize(function() {    
	    $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	    $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	    updateGraph();
	});
});
