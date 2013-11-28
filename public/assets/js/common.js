function parseChartData(data, x, y) {
	var DataPool = [];
	for (var i = 0; i < data.length; i++) {
		var Unit = data[i];
		DataPool.push([ parseFloat(Unit[x]) , parseFloat(Unit[y]) ]);
	}
	return DataPool;
}

function saveUserDefaults(guid, x, y) {
	$.ajax({
        type: "POST",
        //the url where you want to sent the userName and password to
        url: '/api/setdefaults/' + guid,
        //dataType: 'json',
        async: true,
        //json object to sent to the authentication url
        data: "data="+JSON.stringify({ x: x, y: y }),
        success: function (resp) {
       		console.log(resp.result);
        },
        error: function (err) {
        	console.log(err);
        }
    });
}

function getUserDefaults(guid, x, y) {
	$.getJSON( "/api/getdefaults/" + guid, function( data ) {
		$(x).val(data.x);
		$(y).val(data.y);
	});
}