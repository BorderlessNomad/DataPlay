'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

function updateGraph() {
    var graphData = parseChartData(window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()),
        graphAxes = {
            x: $("#pickxaxis").val(),
            y: $("#pickyaxis").val()
        };
    if (!DataCon.graph) {
        $("#chart").html('');
        DataCon.graph = new PGLinesChart("#chart", null, graphData, graphAxes, null);
    } else {
        DataCon.graph.updateChart(graphData, graphAxes);
    }
}

function LightUpBookmarks() {
    // Do nothing
}

window.Annotations = [];

function SavePoint(x, y) {
    var SaveObj = {
        xax: $("#pickxaxis").val(),
        yax: $("#pickyaxis").val(),
        x: x,
        y: y,
        guid: guid
    };
    Annotations.push(SaveObj);
    // Save it as well and rewrite the titlebar URL to be the shareable one.
    $.ajax({
        type: "POST",
        url: '/api/setbookmark/',
        data: "data=" + JSON.stringify(Annotations),
        success: function(resp) {
            window.history.pushState('page2', 'Title', '/viewbookmark/' + resp);
            //Saved and sound
        },
        error: function(err) {
            console.log(err);
        }
    });
}

$(document).ready(function() {

    $("#placeholder").height($(window).height() * 0.8).width($(window).width() * 0.6);
    $(".wikidata").height($(window).height() * 0.8).width($(window).width() * 0.2);
    $('#BubbleLink').click(function() {
        location.href = "/bubble/" + guid;
    });
    $.getJSON("/api/getinfo/" + guid, function(data) {
        $('#FillInDataSet').html(data.Title);
        $("#wikidata").html(data.Notes);
    });

    function populatedKeys(Keys) {
        $("#pickxaxis").empty();
        $("#pickyaxis").empty();
        for (var i = 0; i < Keys.length; i++) {
            $("#pickxaxis").append($("<option></option>").val(Keys[i]).text(Keys[i]));
        }
        for (var i = 0; i < Keys.length; i++) {
            $("#pickyaxis").append($("<option></option>").val(Keys[i]).text(Keys[i]));
        }
        if (Keys.length > 1) {
            $("#pickyaxis").val(Keys[1]);
        }
        // Get the last user preferences if any ... and update graph
        getUserDefaults(guid, $("#pickxaxis"), $("#pickyaxis"), updateGraph);
    }

    $.getJSON("/api/getdata/" + guid, function(data) {
        //console.log(data);
        // so data is an array of shit.
        window.DataSet = data;
        var Keys = [];
        if (data.length) {
            DataCon.patterns = {}
            for (var key in data[0]) {
                Keys.push(key);
                // Pattern recognition
                DataCon.patterns[key] = getPattern(data[0][key]);
            }
        }
        $.ajax({
            async: false,
            url: "/api/identifydata/" + guid,
            dataType: "json",
            success: function(data) {
                var Cols = data.Cols;
                var Keys = [];

                for (var i = 0; i < Cols.length; i++) {
                    if (Cols[i].Sqltype === "int" || Cols[i].Sqltype === "bigint") {
                        Keys.push(Cols[i].Name);
                    }
                }
                populatedKeys(Keys);
            }
        });
        // populatedKeys(Keys);
    });

    $('#SetupOverlay').on("click", function() {
        // Okay so we need to put into local storage the current GUID so when 
        // the overlay panal is called we can use it later on the overlay page.
        localStorage['overlay1'] = guid;
        window.location.href = '/search/overlay';
    });

    window.ReJigGraph = function() {
        saveUserDefaults(guid, $("#pickxaxis").val(), $("#pickyaxis").val());
        updateGraph();
    };

    $(window).resize(function() {
        DataCon.graph = null;
        updateGraph();
    });

});
