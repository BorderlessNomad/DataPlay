'use strict';

$( document ).ready(function() {
	$.getJSON( "/api/user", function( data ) {
		$('#FillInUserName').text(data.Username);
	});
});
