'use strict';


$( document ).ready(function() {
	setInterval(function() {
		var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
		$.getJSON( "/api/getimportstatus/"+guid, function( data ) {
			if(data.state !== "offline") {
				location.href = "http://" + location.host + "/view/" + guid;
			}
		});
	}, 1000);
});
