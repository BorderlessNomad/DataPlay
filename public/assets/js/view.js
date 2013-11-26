'use strict';

$( document ).ready(function() {
	var HavePopulatedKeys = false;
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
	$.getJSON( "/api/getinfo/" + guid, function( data ) {
		$('#FillInDataSet').html(data.Title);
		$(".wikidata").html(data.Notes);
	});
	function PopulatedKeys (Keys) {
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
		console.log(data);
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
				PopulatedKeys(Keys);
				HavePopulatedKeys = true;
			}
			
			DataPool.push([parseInt(Unit[Keys[0]]),parseInt(Unit[Keys[1]])]);
		}
		console.log(DataPool);
		$.plot("#placeholder", [ DataPool ]);
	});
	window.ReJigGraph = function() {
		var DataPool = [];
		var data = window.DataSet;
		for (var i = 0; i < data.length; i++) {
			var Unit = data[i];
			DataPool.push([ parseInt(Unit[$("#pickxaxis").val()]) , parseInt(Unit[$("#pickyaxis").val()]) ]);
		}

		$.plot("#placeholder", [ DataPool ]);
	};
});
