'use strict';

$( document ).ready(function() {
	$('#SBox').keyup(function() {
		if($('#SBox').val() !== "") {
			$.getJSON( "/api/search/" + $('#SBox').val(), function( data ) {
				$('#ResultsTable').empty();
				$('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
				for (var i = data.length - 1; i >= 0; i--) {
					$('#ResultsTable').append("<tr id=\"" + data[i].GUID + "\"><td>"+data[i].Title+"</td><td>" + data[i].GUID + "</td></tr>");
					$.getJSON( "/api/getquality/" + data[i].GUID, function( data ) {
						if(data.Amount > 2) {
							$('#'+ data.Request).addClass("success");
						} else {
							$('#'+ data.Request).addClass("danger");
						}
					});
					console.log(data);
				}
			});
		}
	});
});
