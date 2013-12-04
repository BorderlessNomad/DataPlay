'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

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
    // var SaveObj = {
    // 	xax: $("#pickxaxis").val(),
    // 	yax: $("#pickyaxis").val(),
    // 	x: x,
    // 	y: y,
    // 	guid: guid
    // };
    // Annotations.push(SaveObj);
    // // Save it as well and rewrite the titlebar URL to be the shareable one.
    // $.ajax({
    // 	type: "POST",
    // 	url: '/api/setbookmark/',
    // 	data: "data=" + JSON.stringify(Annotations),
    // 	success: function(resp) {
    // 		window.history.pushState('page2', 'Title', '/viewbookmark/' + resp);
    // 		//Saved and sound
    // 	},
    // 	error: function(err) {
    // 		console.log(err);
    // 	}
    // });
}

function LightUpBookmarks() {
    // do nothing
}

$(document).ready(function() {
    // First thing we need to prepare the table. I'm going to make a 3x3 grid


    // $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
    // $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);

    $.getJSON("/api/getinfo/" + guid, function(data) {
        $('#FillInDataSet').html(data.Title);
        $(".wikidata").html(data.Notes);
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

    $.getJSON("/api/getreduceddata/" + guid, function(data) {
        //console.log(data);
        // so data is an array of shit.
        window.DataSet = data;
        var Keys = [];
        if (data.length) {
            for (var key in data[0]) {
                Keys.push(key);
            }
        }
        // populatedKeys(Keys);
        var Entropy = Math.pow(2, Keys.length);
        var TableHandle = $('#GridTable')
        var count = 0;
        for (var i = 0; i < Entropy/3; i++) {
            $('#GridTable').append('<tr><td><div style="width: 250px;" id="Cell' + count + '"></div></td><td><div style="width: 250px;" id="Cell' + (count + 1) + '"></div></td><td><div style="width: 250px;" id="Cell' + (count + 2) + '"></div></td></tr>');
            count = count + 3;
        }
        var CellCount = 0;
        var k1 = 0;
        var k2 = 0;
        var DCG = [];
        for (var i = 0; i < 9; i++) {
            console.log(i)
            DCG[i] = new PGChart(
                "#Cell" + i,
                null,
                parseChartData(data, Keys[k1], Keys[k2]), {
                    x: Keys[k1],
                    y: Keys[k2]
                },
                100,
                100,
                1
            );
            k1++;
            if (k1 === k2) {
                k1++;
            }
            if (k1 === Keys.length) {
                k1 = 0;
                k2++;
            }
        }
        window.dcg = DCG;
        d3.selectAll('.axis .tick').remove() // Remove the axis to make it look more clean and less
        // insane.

        //drawGraph(parseChartData(data, Keys[0], Keys[1]));
    });

    // window.ReJigGraph = function() {
    // 	drawGraph(parseChartData( window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()));
    // 	saveUserDefaults(guid, $("#pickxaxis").val(), $("#pickyaxis").val());
    // 	updateGraph();
    // };

});
