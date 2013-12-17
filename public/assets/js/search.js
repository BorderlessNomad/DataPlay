'use strict';

$(document).ready(function() {
    var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
    $('#SBox').keyup(function() {
        if ($('#SBox').val() !== "") {
            $.getJSON("/api/search/" + $('#SBox').val(), function(data) {
                $('#ResultsTable').empty();
                $('#ResultsTable').append("<td><tr>Title</tr><tr>Guid</tr></td>");
                for (var i = data.length - 1; i >= 0; i--) {
                    $('#ResultsTable').append("<tr id=\"" + data[i].GUID + "\"><td>" + data[i].Title + "</td><td id=\"Note" + data[i].GUID + "\">" + data[i].GUID + "</td></tr>");
                    $('#' + data[i].GUID).addClass("success");
                    if (guid !== "overlay")
                        $('#Note' + data[i].GUID).html("<a href=\"/view/" + data[i].GUID + "\">View &raquo;</a>&nbsp;<a href=\"/grid/" + data[i].GUID + "\">Grid &raquo;</a>");
                    else
                        $('#Note' + data[i].GUID).html("<a href=\"/overlay/" + data[i].GUID + "\">Overlay &raquo;</a>");


                    // $.getJSON("/api/getquality/" + data[i].GUID, function(data) {
                    //     if (data.Amount > 2) {
                    //         $.getJSON("/api/getimportstatus/" + data[i].GUID, function(dota) {
                    //             if (dota.State == "offline") {
                    //                 $('#Note' + dota.Request).html("<a>Cannot import in the middle of a overlay</a>");
                    //             } else {

                    //             }
                    //         });
                    //         $('#' + data[i].GUID).addClass("success");
                    //     } else {
                    //         $('#' + data[i].GUID).addClass("danger");
                    //         $('#Note' + data[i].GUID).html("Data Quality Too Poor");
                    //     }
                    // });
                    console.log(data);
                }
            });
        }
    });
});
