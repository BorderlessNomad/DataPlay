'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

function updateGraph() {
	DataCon.graph.updateChart(
		parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), 
		{x: $("#pickxaxis").val(), y: $("#pickyaxis").val()}
	);
}

window.Annotations = [];
function SavePoint (x,y) {
	var SaveObj = {
		xax: $("#pickxaxis").val(),
		yax: $("#pickyaxis").val(),
		x: x,
		y: y,
		guid: guid
	};
	Annotations.push(SaveObj);
	// Save it as well and rewrite the titlebar URL to be the shareable one.
	var responceid = "";
	$.ajax({
		type: "POST",
		url: '/api/setbookmark/',
		data: "data=" + JSON.stringify(Annotations),
		success: function(resp) {
			alert(resp);
		},
		error: function(err) {
			console.log(err);
		}
	});
	
	// Now set the URL bar to the correct URL...

}

$( document ).ready(function() {
	$("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	$(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);

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
		// Get the last user preferences if any ... and update graph
		getUserDefaults(guid, $("#pickxaxis"), $("#pickyaxis"), updateGraph);
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
		DataCon.graph = new PGChart(
			"#placeholder", 
			null,
			parseChartData(data, Keys[0], Keys[1]),
			{x: Keys[0], y: Keys[1]}
		);

		//drawGraph(parseChartData(data, Keys[0], Keys[1]));
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
	    $("#placeholder").html('');
	    DataCon.graph = new PGChart(
			"#placeholder", 
			null,
			parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), 
			{x: $("#pickxaxis").val(), y: $("#pickyaxis").val()}
		);
	});

});
