'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];
guid = guid.split('#')[0];

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

function GetURL() {
    // This is called when the coffee script version wants the data URL
    // "/api/getcsvdata/hips/Hospital/60t69"
    //var url = "/api/getcsvdata/" + guid + "/" + $("#pickxaxis").val() + "/" + $("#pickyaxis").val() + "/";
    var url = "/api/getdata/" + guid;
    return url;
};

function createBubblesChart(data) {
    var dataset = [];
    data.forEach(function(d) {
        dataset.push({name: d[$("#pickxaxis").val()], count: d[$("#pickyaxis").val()]});
    });
    $('#chart').html('');
    DataCon.chart = new PGBubblesChart("#chart", {top: 5, right: 0, bottom: 0, left: 0}, dataset, null, null);
}

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
                    XVars.push(Cols[i].Name);
                    YVars.push(Cols[i].Name);
                }
            }
            populateKeys(XVars, YVars);
            // we are storing the current text in the search component
            //  just to make things easy
            var key = decodeURIComponent(location.search).replace("?","");
            // select the current text in the drop-down
            $("#text-select").val(key);
            // bind change in drop down to change the search url and reset the hash url
            d3.select("#pickxaxis").on("change", function(e) {
                d3.json(GetURL(), createBubblesChart);
            });
            d3.select("#pickyaxis").on("change", function(e) {
                d3.json(GetURL(), createBubblesChart);
            });
            //load our data
            d3.json(GetURL(), createBubblesChart)
        }
    });  
    // var dataset = [
    //     {name: 'pepe', count: '234'},
    //     {name: 'juan', count: '24'},
    //     {name: 'rere', count: '23'},
    //     {name: 'fdfh', count: '34'},
    //     {name: 'rrrr', count: '55'},
    //     {name: 'pepsse', count: '33'},
    //     {name: 'asa', count: '112'},
    //     {name: 'sasa', count: '4'},
    //     {name: 'asasas', count: '12'},
    //     {name: 'dddff', count: '32'},
    //     {name: 'pepe', count: '12'},
    //     {name: 'iii', count: '32'},
    //     {name: 'iiii', count: '44'},
    //     {name: 'jjjj', count: '56'},
    //     {name: 'pppp', count: '88'},
    //     {name: 'qqqqq', count: '76'},
    //     {name: 'aaarrrr', count: '56'},
    //     {name: 'ytuy', count: '9'},
    //     {name: 'sdfbvcb', count: '30'},
    //     {name: 'fdsf', count: '19'}

    // ];
    // $('#chart').html('');
    // DataCon.chart = new PGBubblesChart("#chart", {top: 5, right: 0, bottom: 0, left: 0}, dataset, null, null);
});
