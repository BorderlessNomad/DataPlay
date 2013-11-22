'use strict';

$( document ).ready(function() {
	$('#SBox').keyup(function() {
		if($('#SBox').val() !== "") {
			$.getJSON( "/api/search/" + $('#SBox').val(), function( data ) {
				$('#ResultsTable').empty();
				$('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
				for (var i = data.length - 1; i >= 0; i--) {
					$('#ResultsTable').append("<tr id=\"" + data[i].GUID + "\"><td>"+data[i].Title+"</td><td id=\"Note" + data[i].GUID + "\">" + data[i].GUID + "</td></tr>");
					$.getJSON( "/api/getquality/" + data[i].GUID, function( data ) {
						if(data.Amount > 2) {
							$.getJSON( "/api/getimportstatus/" + data.Request, function( dota ) {
								if(dota.State == "offline") {
									$('#Note'+ dota.Request).html("<a href=\"/import/" + dota.Request + "\">Click here to import</a>");
								} else {
									$('#Note'+ dota.Request).html("<a href=\"/view/" + dota.Request + "\">View &raquo;</a>");
								}
								
							});
							$('#'+ data.Request).addClass("success");
						} else {
							$('#'+ data.Request).addClass("danger");
							$('#Note'+ data.Request).html("Data Quality Too Poor");
						}
					});
					console.log(data);
				}
			});
		}
	});
});
