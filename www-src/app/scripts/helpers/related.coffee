'use strict'

class RelatedCharts
	constructor: (chartsRelated) ->
		@chartsRelated = chartsRelated

	preview: false
	allowed: ['line', 'bar', 'row', 'column', 'pie', 'bubble']
	count: 3
	loading:
		related: false
	offset:
		related: 0
	limit:
		related: false
	max:
		related: 0
	chartsRelated: []

	xTicks: 6
	width: 275
	height: 200
	margin:
		top: 10
		right: 10
		bottom: 30
		left: 70
	marginPreview:
		top: 25
		right: 25
		bottom: 25
		left: 25

	monthNames: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]

	humanDate: (date) =>
		"#{date.getDate()} #{@monthNames[date.getMonth()]}, #{date.getFullYear()}"

	findById: (id) ->
		data = _.where(@chartsRelated,
			id: id
		)

		if data?[0]? then data[0] else null

	setPreview: (bool = false) =>
		@preview = bool

	setCustomMargin: (margin = ({top:0,right:0,bottom:0,left:0})) ->
		@customMargin = margin

	isPlotAllowed: (type) ->
		if type in @allowed then true else false

	hasRelatedCharts: () ->
		Object.keys(@chartsRelated).length

	getXScale: (data) ->
		xScale = switch data.patterns[data.xLabel].valuePattern
			when 'label'
				d3.scale.ordinal()
					.domain data.ordinals
					.rangeBands [0, @width]
			when 'date'
				d3.time.scale()
					.domain d3.extent data.group.all(), (d) -> d.key
					.range [0, @width]
			else
				d3.scale.linear()
					.domain d3.extent data.group.all(), (d) -> parseInt d.key
					.range [0, @width]

		xScale

	getXUnits: (data) ->
		xUnits = switch data.patterns[data.xLabel].valuePattern
			when 'date' then d3.time.years
			when 'intNumber' then dc.units.integers
			when 'label', 'text' then dc.units.ordinal
			else dc.units.ordinal

		xUnits

	getYScale: (data) ->
		yScale = switch data.patterns[data.yLabel].valuePattern
			when 'label'
				d3.scale.ordinal()
					.domain data.ordinals
					.rangeBands [0, @height]
			when 'date'
				d3.time.scale()
					.domain d3.extent data.group.all(), (d) -> d.value
					.range [0, @height]
			else
				d3.scale.linear()
					.domain d3.extent data.group.all(), (d) -> parseInt d.value
					.range [0, @height]
					.nice()

		yScale

	lineChartPostSetup: (chart) =>
		data = @findById chart.anchorName()

		data.entry = crossfilter data.values
		data.dimension = data.entry.dimension (d) -> d.x
		data.group = data.dimension.group().reduceSum (d) -> d.y

		chart.dimension data.dimension
		chart.group data.group

		data.ordinals = []
		data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

		chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

		if @preview
			chart.xAxis().ticks(0).tickFormat (v) -> ""
			chart.yAxis().ticks(0).tickFormat (v) -> ""
			chart.margins @marginPreview
		else
			chart.xAxis().ticks @xTicks

		if @customMargin?
			chart.margins @customMargin

		chart.xAxisLabel false, 0
		chart.yAxisLabel false, 0

		chart.x @getXScale data

		return

	rowChartPostSetup: (chart) =>
		data = @findById chart.anchorName()

		data.entry = crossfilter data.values
		data.dimension = data.entry.dimension (d) -> d.x
		data.group = data.dimension.group().reduceSum (d) -> d.y

		chart.dimension data.dimension
		chart.group data.group

		data.ordinals = []
		data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

		chart.colorAccessor (d, i) -> i + 1

		if @preview
			chart.xAxis().ticks(0).tickFormat (v) -> ""
			chart.yAxis().ticks(0).tickFormat (v) -> ""
			chart.margins @marginPreview
		else
			chart.xAxis().ticks @xTicks

		if @customMargin?
			chart.margins @customMargin

		chart.x @getYScale data

		chart.xUnits @getXUnits data if data.ordinals?.length > 0

		return

	columnChartPostSetup: (chart) =>
		data = @findById chart.anchorName()

		data.entry = crossfilter data.values
		data.dimension = data.entry.dimension (d) -> d.x
		data.group = data.dimension.group().reduceSum (d) -> d.y

		chart.dimension data.dimension
		chart.group data.group

		data.ordinals = []
		data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

		chart.colorAccessor (d, i) -> i + 1

		if @preview
			chart.xAxis().ticks(0).tickFormat (v) -> ""
			chart.yAxis().ticks(0).tickFormat (v) -> ""
			chart.margins @marginPreview
		else
			chart.xAxis().ticks @xTicks

		if @customMargin?
			chart.margins @customMargin

		chart.x @getXScale data

		chart.xUnits @getXUnits data if data.ordinals?.length > 0

		return

	pieChartPostSetup: (chart) =>
		data = @findById chart.anchorName()

		data.entry = crossfilter data.values
		data.dimension = data.entry.dimension (d) =>
			if data.patterns[data.xLabel].valuePattern is 'date'
				return @humanDate d.x
			x = if d.x? and (d.x.length > 0 || data.patterns[data.xLabel].valuePattern is 'date') then d.x else "N/A"
		data.groupSum = 0
		data.group = data.dimension.group().reduceSum (d) ->
			y = Math.abs parseFloat d.y
			data.groupSum += y
			y

		chart.dimension data.dimension
		chart.group data.group

		chart.colorAccessor (d, i) -> i + 1

		chart.renderLabel false

		return

	bubbleChartPostSetup: (chart) =>
		data = @findById chart.anchorName()

		minR = null
		maxR = null

		data.entry = crossfilter data.values
		data.dimension = data.entry.dimension (d) ->
			z = Math.abs parseInt d.z

			if not minR? or minR > z
				minR = if z is 0 then 1 else z

			if not maxR? or maxR <= z
				maxR = if z is 0 then 1 else z

			"#{d.x}|#{d.y}|#{d.z}"

		data.group = data.dimension.group().reduceSum (d) -> d.y

		chart.dimension data.dimension
		chart.group data.group

		data.ordinals = []
		for d in data.group.all() when d not in data.ordinals
			data.ordinals.push d.key.split("|")[0]

		chart.keyAccessor (d) -> d.key.split("|")[0]
		chart.valueAccessor (d) -> d.key.split("|")[1]
		chart.radiusValueAccessor (d) ->
			r = Math.abs parseInt d.key.split("|")[2]
			if r >= minR then r else minR

		chart.x switch data.patterns[data.xLabel].valuePattern
			when 'label'
				d3.scale.ordinal()
					.domain data.ordinals
					.rangeBands [0, @width]
			when 'date'
				d3.time.scale()
					.domain d3.extent data.group.all(), (d) -> d.key.split("|")[0]
					.range [0, @width]
			else
				d3.scale.linear()
					.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[0]
					.range [0, @width]

		chart.y switch data.patterns[data.xLabel].valuePattern
			when 'label'
				d3.scale.ordinal()
					.domain data.ordinals
					.rangeBands [0, @height]
			when 'date'
				d3.time.scale()
					.domain d3.extent data.group.all(), (d) -> d.key.split("|")[1]
					.range [0, @height]
			else
				d3.scale.linear()
					.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[1]
					.range [0, @height]

		rScale = d3.scale.linear()
			.domain d3.extent data.group.all(), (d) -> Math.abs parseInt d.key.split("|")[2]
		chart.r rScale

		if @preview
			chart.xAxis().ticks(0).tickFormat (v) -> ""
			chart.yAxis().ticks(0).tickFormat (v) -> ""
			chart.margins @marginPreview
		else
			chart.xAxis().ticks @xTicks

		if @customMargin?
			chart.margins @customMargin

		# chart.label (d) -> x = d.key.split("|")[0]

		chart.title (d) ->
			x = d.key.split("|")[0]
			y = d.key.split("|")[1]
			z = d.key.split("|")[2]
			"#{data.xLabel}: #{x}\n#{data.yLabel}: #{y}\n#{data.zLabel}: #{z}"

		minRL = Math.log minR
		maxRL = Math.log maxR
		scale = Math.abs Math.log (maxRL - minRL) / (maxR - minR)

		chart.maxBubbleRelativeSize scale / 100

		return

window.RelatedCharts = RelatedCharts
