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

	//console.log(options);
	var myChart = new xChart('line', options, '#placeholder');
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
		//$.plot("#placeholder", [ DataPool ]);
	});

	$.getJSON( "/api/getdefaults/" + guid, function( data ) {
		//console.log(data);
	});

	window.ReJigGraph = function() {
		drawGraph(parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()));
		//$.plot("#placeholder", [ DataPool ]);
	};
	$('#SetupOverlay').on("click",function() {
		// Okay so we need to put into local storage the current GUID
		// so when the overlay panal is called we can use it later on the
		// overlay page.
		localStorage['overlay1'] = guid;

		$.ajax({
	        type: "POST",
	        //the url where you want to sent the userName and password to
	        url: '/api/setdefaults/' + guid,
	        dataType: 'json',
	        async: false,
	        //json object to sent to the authentication url
	        data: JSON.stringify(window.DataSet),
	        success: function (resp) {
	       		 console.log(resp.result);
	        }
	    });

		window.location.href = '/search/overlay';
	});
});
