'use strict';

$( document ).ready(function() {
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
	$.getJSON( "/api/search/" + guid, function( data ) {

	}
});
