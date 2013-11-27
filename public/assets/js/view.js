'use strict';

function drawGraph(data) {
	var options = {
		xScale: "linear",
		yScale: "linear",
		type: "line",
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

	console.log(JSON.stringify(options));

	var myChart = new xChart('line', options, '#placeholder');
}

$( document ).ready(function() {
	var HavePopulatedKeys = false;
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
	}
	$.getJSON( "/api/getdata/" + guid, function( data ) {
		// $('#FillInDataSet').html(data.Title);
		//console.log(data);
		// so data is a array of shit.
		window.DataSet = data;
		var DataPool = [];
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			var Keys = [];
			for(var key in Unit) {
				Keys.push(key);
			}
			if(!HavePopulatedKeys) {
				console.log("boop");
				populatedKeys(Keys);
				HavePopulatedKeys = true;
			}
			
			DataPool.push([parseFloat(Unit[Keys[0]]),parseFloat(Unit[Keys[1]])]);
		}
		console.log(DataPool);

		drawGraph(DataPool);

		//$.plot("#placeholder", [ DataPool ]);
	});
	window.ReJigGraph = function() {
		var DataPool = [];
		var data = window.DataSet;
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			DataPool.push([ parseFloat(Unit[$("#pickxaxis").val()]) , parseFloat(Unit[$("#pickyaxis").val()]) ]);
		}

		drawGraph(DataPool);

		//$.plot("#placeholder", [ DataPool ]);
	};
	$('#SetupOverlay').on("click",function() {
		// Okay so we need to put into local storage the current GUID
		// so when the overlay panal is called we can use it later on the
		// overlay page.
		localStorage['overlay1'] = guid;
		window.location.href = '/search/overlay';
	});
});
