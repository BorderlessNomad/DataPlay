'use strict'

settings =
	margin:
		top: 15
		bottom: 25
		right: 55
		left: 55
	marginPreview:
		top: 0
		bottom: 0
		right: 0
		left: 0
	xTicks: 8

typeDictionary =
	bar: 'column'
	area: 'line'

tickFormatFunc = (type) ->
	(d) ->
		if type is 'none'
			return ''
		if type is 'date'
			return d3.time.format("%d-%m-%Y") new Date d
		return d3.format(",f") d

optionsList =
	line:
		chart:
			type: "multiChart"
			height: 450
			margin: _.cloneDeep settings.margin
			x: (d, i) -> i
			y: (d) -> d[1]
			color: d3.scale.category10().range()
			transitionDuration: 0
			xAxis:
				axisLabel: ""
				showMaxMin: false
				tickFormat: tickFormatFunc()
				ticks: settings.xTicks
			yAxis1:
				orient: 'left'
				axisLabel: ""
				tickFormat: tickFormatFunc()
				showMaxMin: false
				highlightZero: false
			yAxis2:
				orient: 'right'
				axisLabel: ""
				tickFormat: tickFormatFunc()
				showMaxMin: false
				highlightZero: false
			yDomain1: [0, 1000]
			yDomain2: [0, 1000]
	column:
		chart:
			type: "multiChart"
			height: 450
			margin: _.cloneDeep settings.margin
			x: (d, i) -> i
			y: (d) -> d[1]
			color: d3.scale.category10().range()
			transitionDuration: 0
			xAxis:
				axisLabel: ""
				showMaxMin: false
				tickFormat: tickFormatFunc()
				ticks: settings.xTicks
			yAxis1:
				orient: 'left'
				axisLabel: ""
				tickFormat: tickFormatFunc()
				showMaxMin: false
				highlightZero: false
			yAxis2:
				orient: 'right'
				axisLabel: ""
				tickFormat: tickFormatFunc()
				showMaxMin: false
				highlightZero: false
			bars1:
				groupSpacing: 0.5
			bars2:
				groupSpacing: 0.5
			yDomain1: [0, 1000]
			yDomain2: [0, 1000]
	scatter:
		chart:
			type: "scatterChart"
			height: 450
			margin: _.cloneDeep settings.margin
			color: d3.scale.category10().range()
			scatter:
				onlyCircles: true
			transitionDuration: 0
			xAxis:
				axisLabel: ""
				showMaxMin: false
				tickFormat: tickFormatFunc()
				ticks: settings.xTicks
			yAxis:
				axisLabel: ""
				showMaxMin: false
				tickFormat: tickFormatFunc()
				highlightZero: false
			yAxis1: {}
			yAxis2: {}
			yDomain1: [0, 1000]
			yDomain2: [0, 1000]

class CorrelatedChart
	constructor: (type, data = []) ->
		if type? then @generate type, data

	preview: false
	type: ''
	options: {}
	data: []
	error: null
	info: {}

	generate: (type, data = []) =>
		if type? and typeDictionary[type]? then type = typeDictionary[type]
		if type? and optionsList[type]?
			@type = type
			@options = _.cloneDeep optionsList[type]
			@data = data
		else
			@error = 'Type not supported'

	setAxisTypes: (xAxis, yAxis1, yAxis2) =>
		items = {xAxis: xAxis, yAxis1: yAxis1, yAxis2: yAxis2}

		if @options?.chart?
			Object.keys(items).forEach (key) =>
				if items[key]
					@options.chart[key].tickFormat = tickFormatFunc items[key]

	setPreview: (bool = false) =>
		@preview = bool
		if bool?
			margin = settings.marginPreview
		else
			margin = settings.margin

		if @options?.chart?.margin?
			if top? then @options.chart.margin.top = margin.top
			if bottom? then @options.chart.margin.bottom = margin.bottom
			if left? then @options.chart.margin.left = margin.left
			if right? then @options.chart.margin.right = margin.right

	setSize: (width, height) =>
		if @options?.chart?
			if width? then @options.chart.width = width
			if height? then @options.chart.height = height

	setMargin: (top, bottom, left, right) =>
		if @options?.chart?.margin?
			if top? then @options.chart.margin.top = top
			if bottom? then @options.chart.margin.bottom = bottom
			if left? then @options.chart.margin.left = left
			if right? then @options.chart.margin.right = right

	setLegend: (flag) =>
		if @options?.chart? and flag? and typeof flag is 'boolean'
			@options.chart.showLegend = flag

	setTooltips: (flag) =>
		if @options?.chart? and flag? and typeof flag is 'boolean'
			@options.chart.tooltips = flag


	translateData: (values, type) =>
		normalise = (d) ->
			if typeof d is 'string'
				if not isNaN Date.parse d
					return Date.parse d
				if not isNaN parseFloat d
					return parseFloat d
			return d

		values.map (v) ->
			newV =
				x: normalise v.x || 0
				y: parseFloat v.y || 0
			if type is 'scatter'
				newV.size = 3
				newV.shape = 'circle'
			newV

window.CorrelatedChart = CorrelatedChart
