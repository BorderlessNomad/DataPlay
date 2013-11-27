'use strict';

$( document ).ready(function() {
	var HavePopulatedKeys = false;
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
	}
	var Keys2Populated = false;
	$.getJSON( "/api/getdata/" + guid, function( data ) {
		var DataPool2 = [];
		var Keys2 = [];
		console.log(data);
		window.DataSet = data;
		var DataPool = [];
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			var Keys = [];
			for(var key in Unit) {
				Keys.push(key);
			}

			
			DataPool.push([parseFloat(Unit[Keys[0]]),parseFloat(Unit[Keys[1]])]);
		}
		console.log(DataPool);

		$.getJSON( "/api/getdata/" + localStorage['overlay1'], function( dataa ) {
			window.DataSet2 = dataa;
			for (var i = 0; i < dataa.length; i++) {
				var Unit = dataa[i];
				if(!Keys2Populated) {
					for(var key in Unit) {
						Keys2.push(key);
					}
					Keys2Populated = true;
				}
				if(!HavePopulatedKeys && Keys2Populated) {
					console.log("boop");
					populatedKeys(Keys,Keys2);
					HavePopulatedKeys = true;
				}
				DataPool2.push([parseFloat(Unit[Keys2[0]]),parseFloat(Unit[Keys2[1]])]);
			}
		});
		
		$.plot("#placeholder", [ DataPool,DataPool2 ]);
	});
	window.ReJigGraph = function() {
		var DataPool = [];
		var data = window.DataSet;
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			DataPool.push([ parseFloat(Unit[$("#pickx1axis").val()]) , parseFloat(Unit[$("#picky1axis").val()]) ]);
		}
		var DataPool2 = [];
		var data = window.DataSet2;
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			DataPool2.push([ parseFloat(Unit[$("#pickx2axis").val()]) , parseFloat(Unit[$("#picky2axis").val()]) ]);
		}


		$.plot("#placeholder", [ 
			{ data: DataPool, label: $("#picky1axis").val() },
			{ data: DataPool2, label: $("#picky2axis").val(), yaxis: 2 }
		]);
	};
});
