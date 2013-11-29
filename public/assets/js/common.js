function parseChartData(data, x, y) {
	var DataPool = [];
	for (var i = 0; i < data.length; i++) {
		DataPool.push([ parseFloat(data[i][x]) , parseFloat(data[i][y]) ]);
	}
	return DataPool;
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