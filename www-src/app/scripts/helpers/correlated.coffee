'use strict'

settings =
	margin:
		top: 15
		bottom: 25
		right: 55
		left: 55
	xTicks: 8

typeDictionary =
	bar: 'column'
	area: 'line'

tickFormatFunc = (type) ->
	(d) ->
		if type is 'date'
			return d3.time.format("%d-%m-%Y") new Date d
		return d3.format(",f") d

optionsList =
	line:
		chart:
			type: "multiChart"
			height: 450
			margin: settings.margin
			x: (d, i) -> i
			y: (d) -> d[1]
			color: d3.scale.category10().range()
			transitionDuration: 250
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
			margin: settings.margin
			x: (d, i) -> i
			y: (d) -> d[1]
			color: d3.scale.category10().range()
			transitionDuration: 250
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
			margin: settings.margin
			color: d3.scale.category10().range()
			scatter:
				onlyCircles: true
			transitionDuration: 250
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

class CorrelatedChart
	constructor: (type, data = []) ->
		if type? then @generate type, data

	type: ''
	options: {}
	data: []
	error: null

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

window.CorrelatedChart = CorrelatedChart
