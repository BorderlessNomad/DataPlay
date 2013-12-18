'use strict';

$( document ).ready(function() {
	$.getJSON( "/api/user", function( data ) {
		$('#FillInUserName').text(data.Username);
	});
	$.getJSON( "/api/visited", function( data ) {
		if(data.length != 0) {
			$('#History').empty();
			$('#History').append("<p>Welcome Back, You where last viewing:</p>");
		}
		for (var i = 0; i < data.length; i++) {
			if(data[i][1].length > 85) {
				data[i][1] = data[i][1].substring(0,83) + "...";
			}
			$('#History').append("<a href=\"/view/" + data[i][0] + "\"> "+data[i][1]+" </a></br>");
		}
	});
});
