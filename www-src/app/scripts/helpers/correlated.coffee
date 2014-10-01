'use strict'

settings =
	margin:
		top: 15
		bottom: 25
		right: 55
		left: 55
	xTicks: 8

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
				tickFormat: (d) -> d3.time.format("%d-%m-%Y") new Date d
				ticks: settings.xTicks
			yAxis1:
				orient: 'left'
				axisLabel: ""
				tickFormat: (d) -> d3.format(",f") d
				showMaxMin: false
				highlightZero: false
			yAxis2:
				orient: 'right'
				axisLabel: ""
				tickFormat: (d) -> d3.format(",f") d
				showMaxMin: false
				highlightZero: false
			areas:
				dispatch:
					elementClick: (e) ->
						console.log 'areas elementClick', e
			lines:
				dispatch:
					elementClick: (e) ->
						console.log 'lines elementClick', e
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
				tickFormat: (d) -> d3.time.format("%d-%m-%Y") new Date d
				ticks: settings.xTicks
			yAxis1:
				orient: 'left'
				axisLabel: ""
				tickFormat: (d) -> d3.format(",f") d
				showMaxMin: false
				highlightZero: false
			yAxis2:
				orient: 'right'
				axisLabel: ""
				tickFormat: (d) -> d3.format(",f") d
				showMaxMin: false
				highlightZero: false
			areas:
				dispatch:
					elementClick: (e) ->
						console.log 'areas elementClick', e
			lines:
				dispatch:
					elementClick: (e) ->
						console.log 'lines elementClick', e
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
		if type? and optionsList[type]?
			@type = type
			@options = _.cloneDeep optionsList[type]
			@data = data
		else
			@error = 'Type not supported'


window.CorrelatedChart = CorrelatedChart
