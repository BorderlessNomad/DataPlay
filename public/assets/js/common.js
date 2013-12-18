_.templateSettings.variable = "data"

function getPattern(data) {
    //console.log('--------- Processing pattern for ' + data + ' --------');
    var patterns = {
        date: /^\d{2,4}-\d{1,2}-\d{1,2}$/i,
        label: /[a-z]+/i,
        intNumber: /^\d+$/i,
        floatNumber: /^\d*\.|,\d+$/i,
        percentage: /^\d+(\.|,\d)*%$/i
    };
    for (var pt in patterns) {     
        if (patterns.hasOwnProperty(pt)) {
            //console.log(pt + ' pattern?');
            var isPattern = patterns[pt].exec(data);
            //console.log(isPattern ? isPattern : 'No matches');
            if (isPattern) return pt;               
        }
    }
    return null;
}

function parseDate(input) {
    var parts = input.split('-');
    return new Date(parts[0], parts[1]-1, parts[2]);
}

function formatDate(input) {
    return input.getFullYear()+'-'+(input.getMonth()+1)+'-'+input.getDate();
}

function parseAxisData(data, key) {
    var datapool = [];
    DataCon.patterns || (DataCon.patterns = {});
    for (var i = 0; i < data.length; i++) {
        switch (DataCon.patterns[key]) {
            case 'date': 
                datapool.push(parseDate(data[i][key]));
                break;
            case 'label': 
                datapool.push(data[i][key]);
                break;
            case 'intNumber': 
                datapool.push(parseInt(data[i][key]));
                break;
            default: 
                datapool.push(parseFloat(data[i][key]));
        }
    }
    return datapool;
}

function parseChartData(data, x, y) {
    var xData = parseAxisData(data, x),
                    yData = parseAxisData(data, y),
                    datapool = [];
    for (var i = 0; i < data.length; i++) {
       datapool.push([ xData[i] , yData[i]]);
    }
    return datapool;
}

function quicksort(dataset) {
     if (dataset.length <= 1) {
         return dataset;
     } else {
        var pivot = Math.round(dataset.length/2),
              less = [],
              greater = [];
         for (var i=0; i< dataset.length; i++) {
            if (i != pivot) {
                 dataset[i][0] <= dataset[pivot][0] ? less.push(dataset[i]) : greater.push(dataset[i]);
            }      
         }
         return quicksort(less).concat([dataset[pivot]]).concat(quicksort(greater));
    }
}

function saveUserDefaults(guid, x, y) {
    $.ajax({
        type: "POST",
        url: '/api/setdefaults/' + guid,
        async: true,
        data: "data="+JSON.stringify({ x: x, y: y }),
        success: function (resp) {
            //console.log(resp.result);
        },
        error: function (err) {
            console.log(err);
        }
    });
}

function getUserDefaults(guid, x, y, cb) {
    $.getJSON( "/api/getdefaults/" + guid, function( data ) {
        $(x).val(data.x);
        $(y).val(data.y);
        cb && cb();
    });
}
