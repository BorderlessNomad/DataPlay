'use strict';

window.DataCon || (window.DataCon = {})
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

window.Annotations = [];

function SavePoint(x, y) {

}

function LightUpBookmarks() {
    // do nothing
}

function ProabblyANumber(name, identdata) {
    var RealArray = identdata.Cols;
    for (var i = 0; i < RealArray.length; i++) {
        if(RealArray[i].Name === name && (RealArray[i].Sqltype === "float" || RealArray[i].Sqltype === "bigint" || RealArray[i].Sqltype === "int")) {
            return true;
        }
    }
    return false;
}

$(document).ready(function() {
    // First thing we need to prepare the table. I'm going to make a 3x3 grid

    window.SetAndGo = function(x, y) {
        $.ajax({
            type: "POST",
            url: '/api/setdefaults/' + guid,
            async: true,
            data: "data=" + JSON.stringify({
                x: x,
                y: y
            }),
            success: function(resp) {
                window.location.href = "/view/" + guid;
            },
            error: function(err) {
                console.log(err);
            }
        });
    };

    // $("#placeholder").height($(window).height()*0.8).width($(window).width()*0.6);
    // $(".wikidata").height($(window).height()*0.8).width($(window).width()*0.2);

    $.getJSON("/api/getinfo/" + guid, function(data) {
        $('#FillInDataSet').html(data.Title);
        $("#wikidata").html(data.Notes);
    });


    $.getJSON("/api/getreduceddata/" + guid, function(data) {
        $.getJSON("/api/identifydata/" + guid, function(tabledata) {

            // because of what the grid was meant to show, I am going to make sure it can only plot
            // number to number
            window.DataSet = data;
            var Keys = [];
            DataCon.patterns = {};
            if (data.length) {
                for (var key in data[0]) {
                    Keys.push(key);
                    // Pattern recognition
                    DataCon.patterns[key] = getPattern(data[0][key]);
                }
            }
            var Entropy = Math.pow(2, Keys.length - 1);
            var GridMembers = [];
            var k1 = 0;
            var k2 = 0;
            for (var i = 0; i < Entropy; i++) {
                if(ProabblyANumber(Keys[k2],tabledata)) {
                    GridMembers.push({
                        x: Keys[k1],
                        y: Keys[k2]
                    });
                }
                k1++;
                if (k1 === k2) {
                    k1++;
                }
                if (k1 === Keys.length) {
                    k1 = 0;
                    k2++;
                }
            }
            console.log(GridMembers)

            var TableHandle = $('#GridTable');
            var count = 0;
            for (var i = 0; i < GridMembers.length / 3; i++) {
                $('#GridTable').append(
                    '<tr>' +
                    '<td><div class="gridCell" id="Cell' + count + '"></div></td>' +
                    '<td><div class="gridCell" id="Cell' + (count + 1) + '"></div></td>' +
                    '<td><div class="gridCell" id="Cell' + (count + 2) + '"></div></td>' +
                    '</tr>');
                count = count + 3;
            }
            var CellCount = 0;
            var DCG = [];
            for (var i = 0; i < GridMembers.length; i++) {
                DCG[i] = new PGLinesChart(
                    "#Cell" + i, {
                        top: 5,
                        right: 5,
                        bottom: 5,
                        left: 5
                    },
                    parseChartData(data, GridMembers[i].x, GridMembers[i].y), {
                        x: GridMembers[i].x,
                        y: GridMembers[i].y
                    },
                    null
                );
                $('#Cell' + i).parent().append('<a onclick="SetAndGo(\'' + Keys[k1] + '\',\'' + Keys[k2] + '\');">View</a>');
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
            d3.selectAll('.axis .tick').remove();
            // Remove the axis to make it look more clean and less insane.

        });
    });

});
