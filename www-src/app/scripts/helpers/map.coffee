class MapGenerator
	constructor: (selector = '', data = []) ->
		@selector = selector
		@data = data

		$.getJSON '../regions.json', (data) ->
			MapGenerator.prototype.boundaryPaths = data

	selector: ''
	data: []

	width: 156
	height: 168

	maxvalue: 50

	generate: (data) =>
		if data then @data = data else data = @data

		el = d3.select @selector
		svg = el.append 'svg'
			.attr 'width', @width
			.attr 'height', @height
		trans = svg.append 'g'
			.attr 'transform', "translate(#{@width / 2}, #{@height / 2})"

		data.forEach (c) =>
			@appendCounty trans, c

	appendCounty: (container, county) =>
		displayName = @displayNameDictionary[county.name.toLowerCase()]
		if not displayName?
			displayName = county.name.toLowerCase().replace(/^(ire)/g, '')
			displayName = displayName.substring(0,1).toUpperCase() + displayName.substring(1)

		container.append 'g'
			.attr 'class', 'county'
			.attr 'id', "county-#{county.name}"
			.attr 'data-display', displayName
			.append 'path'
				.attr 'fill', @getColor county
				.attr 'd', @boundaryPaths[county.name]
			.append 'title'
				.html "#{displayName}: #{county.value}"

	getColor: (county, max) =>
		value = if typeof county is 'object' then (county.value || 0) else (county || 0)

		value = value / (max or @maxvalue)

		if typeof county is 'object' and county.name.substring(0, 3) is 'ire' then return '#CCCCCC'
		if value is 0 then return '#FFE455'

		start = { r: 255, g: 228, b: 85 }
		end   = { r: 186, g: 63 , b: 60 }
		rgb = {}

		Object.keys(start).forEach (col) ->
			rgb[col] = Math.round start[col] + ((end[col] - start[col]) * value)

		return '#' + rgb.r.toString(16) + rgb.g.toString(16) + rgb.b.toString(16)

	highlight: (corr) =>
		d3.select "##{corr} path"
			.style 'fill-opacity', 0.8

	unhighlight: (corr) =>
		d3.select "##{corr} path"
			.style 'fill-opacity', null


	locationDictionary:
		'london': 'greaterlondon'
		'manchester': 'greatermanchester'

	displayNameDictionary:
		'dumfriesandgalloway':  'Dumfries and Galloway'
		'eastsussex':           'East Sussex'
		'greaterlondon':        'London'
		'greatermanchester':    'Manchester'
		'herefordandworcester': 'Hereford and Worcester'
		'isleofwight':          'Isle of Wight'
		'midglamorgan':         'Mid Glamorgan'
		'northernireland':      'Northern Ireland'
		'northyorkshire':       'North Yorkshire'
		'orkneyislands':        'Orkney Islands'
		'scottishborders':      'Scottish Borders'
		'southglamorgan':       'South Glamorgan'
		'southyorkshire':       'South Yorkshire'
		'tyneandwear':          'Tyne and Wear'
		'westernisles':         'Western Isles'
		'westglamorgan':        'West Glamorgan'
		'westmidlands':         'West Midlands'
		'westsussex':           'West Sussex'
		'westyorkshire':        'West Yorkshire'

	boundaryPaths: {}


window.MapGenerator = MapGenerator
