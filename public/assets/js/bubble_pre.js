'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

// function updateGraph() {
// 	DataCon.graph.updateChart(
// 		parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), 
// 		{x: $("#pickxaxis").val(), y: $("#pickyaxis").val()}
// 	);
// }


$( document ).ready(function() {

	// http://localhost:3000/api/identifydata/hips
	$.getJSON("/api/identifydata/" + guid,function(data) {
		// {
		//     "Cols": [
		//         {
		//             "Name": "`Hospital`",
		//             "Sqltype": "varchar"
		//         }
		//     ],
		//     "Request": "hips"
		// }
		var Cols = data.Cols;
		var XVars = [];
		var YVars = [];

		for (var i = 0; i < Cols.length; i++) {
			if(Cols[i].Sqltype === "varchar") {
				XVars.push(Cols[i].Name);
			} else {
				YVars.push(Cols[i].Name);
			}
		}
		populateKeys(XVars,YVars);
	});

	// $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	// $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	// $('#BubbleLink').click(function() {
	// 	location.href="/bubble/"+guid;
	// });
	// $.getJSON( "/api/getinfo/" + guid, function( data ) {
	// 	$('#FillInDataSet').html(data.Title);
	// 	$(".wikidata").html(data.Notes);
	// });

	function populateKeys (XVars,YVars) {
		$("#pickxaxis").empty();
		$("#pickyaxis").empty();
		for (var i = 0; i < XVars.length; i++) {
			$("#pickxaxis").append($("<option></option>").val(XVars[i]).text(XVars[i]));
		}
		for (var i = 0; i < YVars.length; i++) {
			$("#pickyaxis").append($("<option></option>").val(YVars[i]).text(YVars[i]));
		}
	}

	// $.getJSON( "/api/getdata/" + guid, function( data ) {
	// 	//console.log(data);
	// 	// so data is an array of shit.
	// 	window.DataSet = data;
	// 	var Keys = [];
	// 	if (data.length) {
	// 		for(var key in data[0]) {
	// 			Keys.push(key);
	// 		}
	// 	}
	// 	populatedKeys(Keys);
	// 	DataCon.graph = new PGChart(
	// 		"#placeholder", 
	// 		null,
	// 		parseChartData(data, Keys[0], Keys[1]),
	// 		{x: Keys[0], y: Keys[1]}
	// 	);

	// 	//drawGraph(parseChartData(data, Keys[0], Keys[1]));
	// });

	// window.ReJigGraph = function() {
	// 	DataCon.graph.updateChart(parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()))
	// 	// drawGraph();
	// 	saveUserDefaults(guid, $("#pickxaxis").val(), $("#pickyaxis").val());
	// 	updateGraph();
	// };

	// $('#SetupOverlay').on("click",function() {
	// 	// Okay so we need to put into local storage the current GUID
	// 	// so when the overlay panal is called we can use it later on the
	// 	// overlay page.
	// 	localStorage['overlay1'] = guid;
	// 	window.location.href = '/search/overlay';
	// });

	// $(window).resize(function() {    
	//     $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
	//     $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);
	//     $("#placeholder").html('');
	//     DataCon.graph = new PGChart(
	// 		"#placeholder", 
	// 		null,
	// 		parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), 
	// 		{x: $("#pickxaxis").val(), y: $("#pickyaxis").val()}
	// 	);
	// });

});
