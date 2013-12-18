'use strict';

window.DataCon || (window.DataCon = {})
DataCon.currChartType = 'enclosure';
var guid = window.location.href.split('/')[window.location.href.split('/').length - 1];

function go2HierarchyChart() {
    var count = _.keys(window.DataSet[0]).length,
        selectorTemplate = _.template($('#selectorTemplate').html()),
        selectorButton = $('<a class="btn btn-lg btn-success" href="javascript:void(0)" role="button">Enclosure Chart</a>'),
        data = { type: 'Key', keys: [] },
        dataset = { key: '', values: null };             
    for (var node in window.DataSet[0]) {
        data.keys.push(node);
    }
    $('#selectors').html('');
    for (var i=0; i<count-1; i++) {
        $('#selectors').append(selectorTemplate(data));
    }
    data.type = 'Value';
    $('#selectors').append(selectorTemplate(data));
    $('#controls').html('').append(selectorButton);
    selectorButton.click(function() {
        $("#chart").html('');
        var keySelectors = $('#selectors select.Key'),
            valueSelector = $('#selectors select.Value'),
            nested = d3.nest();
        for (var i=0; i<keySelectors.length; i++) {
            var val = $(keySelectors.get(i)).val();
            if (val && val!='0') {
                (function(theNest, theNode) {
                    theNest.key(function(d) {
                        return d[theNode];
                    });
                })(nested, $(keySelectors.get(i)).val());
            }                     
        }
        nested = nested.entries(window.DataSet);
        dataset.values = nested;
        console.log(valueSelector.val());
        switch (DataCon.currChartType) {
            case 'enclosure':
                DataCon.chart = new PGEnclosureChart("#chart", {top: 0, right: 0, bottom: 0, left: 0},
                    dataset, null, null, valueSelector.val());
                break;
        }
    });

}

function go2Chart(type) {
    if (DataCon.currChartType != type) {
        DataCon.currChartType = type;
        DataCon.chart = null;
        $('#modes').html('');
        updateChart();
    }
}

function updateChart() {
    var chartData = parseChartData(window.DataSet, $("#pickxaxis").val(), $("#pickyaxis").val()),
        chartAxes = {
            x: $("#pickxaxis").val(),
            y: $("#pickyaxis").val()
        };
    if (!DataCon.chart) {
        $("#chart").html('');
        switch (DataCon.currChartType) {
            case 'bars':
                DataCon.chart = new PGBarsChart("#chart", null, chartData, chartAxes, 30);
                break;
            case 'area':
                DataCon.chart = new PGAreasChart("#chart", null, chartData, chartAxes, null);
                break;
            case 'pie':
                DataCon.chart = new PGPieChart("#chart", null, chartData, null, null);
                break;
            case 'bubbles':
                var dataset = [];
                chartData.forEach(function(d) {
                    dataset.push({name: d[0], count: d[1]});
                });
                DataCon.chart = new PGBubblesChart("#chart", {top: 5, right: 0, bottom: 0, left: 0}, 
                    dataset, null, null);
                break;
            case 'enclosure':
            case 'tree':
            case 'treemap':
                go2HierarchyChart();
                break;
            default:
                DataCon.chart = new PGLinesChart("#chart", null, chartData, chartAxes, null);
        }
        
    } else {
        switch (DataCon.currChartType) {
            case 'enclosure':
                // Special treatment????
                break;
            default: DataCon.chart.updateChart(chartData, chartAxes);
        }
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
        getUserDefaults(guid, $("#pickxaxis"), $("#pickyaxis"), updateChart);
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
                //var Keys = [];
                for (var i = 0; i < Cols.length; i++) {
                    // if (Cols[i].Sqltype === "int" || Cols[i].Sqltype === "bigint" || Cols[i].Sqltype === "float") {
                    //     Keys.push(Cols[i].Name);
                    // }
                    switch (Cols[i].Sqltype) {
                        case "int", "bigint":
                            DataCon.patterns[Cols[i].Name] = 'intNumber';
                        case "float":
                            DataCon.patterns[Cols[i].Name] = 'floatNumber';
                        case 'varchar':
                            // leave pattern as it was recognised by frontend
                        default:
                            // leave pattern as it was recognised by frontend
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
        updateChart();
    };

    $(window).resize(function() {
        DataCon.chart = null;
        updateChart();
    });

});
