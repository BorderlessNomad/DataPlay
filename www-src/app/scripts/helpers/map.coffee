class MapGenerator
	constructor: (selector = '', data = []) ->
		@selector = selector
		@data = data

		$.getJSON '../regions.json', (data) ->
			MapGenerator.prototype.boundaryPaths = data

	selector: ''
	data: []

	width: 180
	height: 148

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
			c.slug = c.name.toLowerCase().replace(/_|\-|\'|\s/g, '');
			@appendRegion trans, c

	appendRegion: (container, region) =>
		container.append 'g'
			.attr 'class', 'region'
			.attr 'id', "region-#{region.slug}"
			.attr 'data-display', region.name
			.append 'path'
				.attr 'fill', @getColor region
				.attr 'stroke', '#000000'
				.attr 'stroke-width', '1'
				.attr 'stroke-opacity', '0.05'
				.attr 'd', @boundaryPaths[region.name]
			.append 'title'
				.html "#{region.name}: #{region.value}"

	getColor: (region, max) =>
		value = if typeof region is 'object' then (region.value || 0) else (region || 0)

		value = value / (max or @maxvalue)

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

	boundaryPaths: {}


window.MapGenerator = MapGenerator
