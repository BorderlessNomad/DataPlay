'use strict';

window.DataCon || (window.DataCon = {})


function updateGraph() {
    DataCon.graph.updateChart(
        parseChartData(window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), {
            x: $("#pickxaxis").val(),
            y: $("#pickyaxis").val()
        }
    );
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
    var responceid = "";
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

    window.ReJigGraph = function() {
        drawGraph(parseChartData(window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()));
        updateGraph();
    };

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

    // Okay so its now critically importent that we graph and put the GUID into
    // Window.guid from the bookmark.

    $.getJSON("/api/getbookmark/" + window.location.href.split('/')[window.location.href.split('/').length - 1], function(bookmarks) {
        window.guid = bookmarks[0].guid;
        $.getJSON("/api/getinfo/" + guid, function(data) {
            $('#FillInDataSet').html(data.Title);
            $(".wikidata").html(data.Notes);
        });

        $.getJSON("/api/getdata/" + guid, function(data) {
            //console.log(data);
            // so data is an array of shit.
            window.DataSet = data;
            var Keys = [];
            if (data.length) {
                for (var key in data[0]) {
                    Keys.push(key);
                }
            }
            populatedKeys(Keys);
            DataCon.graph = new PGChart(
                "#placeholder",
                null,
                parseChartData(data, Keys[0], Keys[1]), {
                    x: Keys[0],
                    y: Keys[1]
                }
            );
            //drawGraph(parseChartData(data, Keys[0], Keys[1]));
        });
        window.Annotations = bookmarks;
        window.LightUpBookmarks = function (argument) {
            window.Annotations = bookmarks;
            var WhatInTheHell = d3.select(window.GraphOBJ.data()).filter;
            for (var i = bookmarks.length - 1; i >= 0; i--) {
                var bm = bookmarks[i];
                GraphOBJ.data().filterr = WhatInTheHell;
                GraphOBJ.data().filterr(function(datum) {
                    return datum[0] === bm.x && datum[1] === bm.y || Math.random() > 0.7;
                }).attr("fill","#DEADBE");

                // WhatInTheHell
            }
        };
        
        
    });





    $(window).resize(function() {
        $("#placeholder").height($(window).height() * 0.8).width($(window).width() * 0.6);
        $(".wikidata").height($(window).height() * 0.8).width($(window).width() * 0.2);
        $("#placeholder").html('');
        DataCon.graph = new PGChart(
            "#placeholder",
            null,
            parseChartData(window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()), {
                x: $("#pickxaxis").val(),
                y: $("#pickyaxis").val()
            }
        );
    });

});
