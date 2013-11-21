'use strict';

$( document ).ready(function() {
	$('#SBox').keydown(function() {
		$('#ResultsTable').empty();
		if($('#SBox').val() !== "") {
			$.getJSON( "/api/search/" + $('#SBox').val(), function( data ) {
				console.log(data);
				$('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
				for (var i = data.length - 1; i >= 0; i--) {
					$('#ResultsTable').append("<tr><td>"+data[i].Title+"</td><td>" + data[i].GUID + "</td></tr>");
					console.log(data[i].Title);
				}
			});
		}
	});
});
