'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
guid = guid.split('#')[0];


$(document).ready(function() {

    $.ajax({
        async: false,
        url: "/api/identifydata/" + guid,
        dataType: "json",
        success: function(data) {
            var Cols = data.Cols;
            var XVars = [];
            var YVars = [];

            for (var i = 0; i < Cols.length; i++) {
                if (Cols[i].Sqltype === "varchar") {
                    XVars.push(Cols[i].Name);
                } else {
                    YVars.push(Cols[i].Name);
                }
            }
            populateKeys(XVars, YVars);
        }
    });

    function populateKeys(XVars, YVars) {
        $("#pickxaxis").empty();
        $("#pickyaxis").empty();
        for (var i = 0; i < XVars.length; i++) {
            $("#pickxaxis").append($("<option></option>").val(XVars[i]).text(XVars[i]));
        }
        for (i = 0; i < YVars.length; i++) {
            $("#pickyaxis").append($("<option></option>").val(YVars[i]).text(YVars[i]));
        }
    }
    window.GetURL = function() {
        // This is called when the coffee script version wants the data URL
        // "/api/getcsvdata/hips/Hospital/60t69"
        var url = "/api/getcsvdata/" + guid + "/" + $("#pickxaxis").val() + "/" + $("#pickyaxis").val() + "/";
        return url;
    };
});
