'use strict';

$( document ).ready(function() {
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
	$.getJSON( "/api/getinfo/" + guid, function( data ) {
		$('#FillInDataSet').html(data.Title);
		$(".wikidata").html(data.Notes);
	});
	$.getJSON( "/api/getdata/" + guid, function( data ) {
		// $('#FillInDataSet').html(data.Title);
		console.log(data);
	});

});
