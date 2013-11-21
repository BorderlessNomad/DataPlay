'use strict';

$( document ).ready(function() {
	$('#SBox').keyup(function() {
		if($('#SBox').val() !== "") {
			$.getJSON( "/api/search/" + $('#SBox').val(), function( data ) {
				$('#ResultsTable').empty();
				$('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
				for (var i = data.length - 1; i >= 0; i--) {
					$('#ResultsTable').append("<tr><td>"+data[i].Title+"</td><td>" + data[i].GUID + "</td></tr>");
					console.log(data[i].Title);
				}
			});
		}
	});
});
