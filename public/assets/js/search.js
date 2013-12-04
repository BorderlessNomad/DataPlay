'use strict';

$( document ).ready(function() {
	var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
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
									$('#Note'+ dota.Request).html("<a>Cannot import in the middle of a overlay</a>");
								} else {
									if(guid !== "overlay")
										$('#Note'+ dota.Request).html("<a href=\"/view/" + dota.Request + "\">View &raquo;</a>&nbsp;<a href=\"/grid/" + dota.Request + "\">Grid &raquo;</a>");
									else
										$('#Note'+ dota.Request).html("<a href=\"/overlay/" + dota.Request + "\">Overlay &raquo;</a>");
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
