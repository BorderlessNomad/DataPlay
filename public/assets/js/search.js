'use strict';

$(document).ready(function() {
    var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
    $('#SBox').keyup(function() {
        if ($('#SBox').val() !== "") {
            $.getJSON("/api/search/" + $('#SBox').val(), function(data) {
                $('#ResultsTable').empty();
                $('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
                for (var i = data.length - 1; i >= 0; i--) {
                	if(data[i].Title != "") {
	                    $('#ResultsTable').append("<tr id=\"" + data[i].GUID + "\"><td style=\"width: 70%;\">" + data[i].Title + "</td><td id=\"Note" + data[i].GUID + "\">" + data[i].GUID + "</td></tr>");
	                    $('#' + data[i].GUID).addClass("success");
	                    if (guid !== "overlay")
	                        $('#Note' + data[i].GUID).html("<a href=\"/view/" + data[i].GUID + "\">View &raquo;</a>&nbsp;<a href=\"/grid/" + data[i].GUID + "\">Grid &raquo;</a>");
	                    else
	                        $('#Note' + data[i].GUID).html("<a href=\"/overlay/" + data[i].GUID + "\">Overlay &raquo;</a>");
	                    console.log(data);
	                }
                }
            });
        }
    });
});
