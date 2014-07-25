'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewCtrl
 # @description
 # # OverviewCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewCtrl', ['$scope', '$routeParams', 'Overview', 'PatternMatcher', ($scope, $routeParams, Overview, PatternMatcher) ->
		$scope.params = $routeParams
		$scope.chartsInfo = []

		$scope.guid = null
		$scope.container = 'body'
		$scope.data = null
		$scope.keys = []
		$scope.cfdata = null
		$scope.dimensions = []
		$scope.dimensionsMap = {}
		$scope.groups = []
		$scope.charts = [
			{ id: 'row', maxEntries: 15 }
			{ id: 'bar', maxEntries: 30 }
			{ id: 'pie', maxEntries: 45 }
			{ id: 'bubble', maxEntries: 60 }
			{ id: 'line' }
		]
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 30

		d3.csv "bower_components/angular-dc/example/stocks/morley.csv", (error, experiments) ->
			ndx = crossfilter experiments
			$scope.runDimension = ndx.dimension (d) ->
				"run-" + d.Run
			$scope.speedSumGroup = $scope.runDimension.group().reduceSum (d) ->
				d.Speed * d.Run

		$scope.getCharts = () ->
			Overview.reducedData $scope.params.id, 10, 100
				.success (data) ->
					if data.length
						patterns = {}
						for key of data[0]
							do (key) ->
								vp = PatternMatcher.getPattern data[0][key]
								kp = PatternMatcher.getKeyPattern key
								patterns[key] = valuePattern: vp, keyPattern: kp
								# Now parse ALL the data
								# TODO get into account key pattern before parsing everything???
								entry[key] = PatternMatcher.parse(entry[key], patterns[key].valuePattern) for entry in data

						console.log "Patterns: ", patterns

						$scope.plot $scope.params.id, {dataset: data, patterns: patterns}, '#charts'

					return
				.error (data) ->
					console.log "Overview::getCharts::Error:", status
					return

			return

		$scope.plot = (guid, data, container, width, height) ->
			$scope.guid = guid
			$scope.data = data
			$scope.container = container

			$scope.width = width if width
			$scope.height = height if height

			$scope.processData()
			$scope.drawCharts()

			return

		$scope.processData = ->
			$scope.keys = (entry for entry of $scope.data.patterns)
			$scope.cfdata = crossfilter $scope.data.dataset#?.slice(0, 10) # TODO: remove this limiting slice

			$scope.dimensions.push($scope.cfdata.dimension (d) -> d[key]) for key in $scope.keys

			if $scope.dimensions.length > 1
				for i in [0..$scope.dimensions.length-2]
					do (i) =>
						for j in [i+1..$scope.dimensions.length-1]
							do (j) =>
								$scope.addGroup i, j
								$scope.addGroup j, i
								return

				#console.log entry.group.all() for entry in $scope.groups

		$scope.getFilteredDataset = ->
			$scope.dimensions[0]?.bottom Infinity

		$scope.addGroup = (i, j) ->
			xKey = $scope.keys[i]
			xKeyPattern = $scope.data.patterns[xKey].keyPattern
			xPattern = $scope.data.patterns[xKey].valuePattern

			yKey = $scope.keys[j]
			yKeyPattern = $scope.data.patterns[yKey].keyPattern
			yPattern = $scope.data.patterns[yKey].valuePattern

			# console.log xKey, xKeyPattern, xPattern, yKey, yKeyPattern, yPattern

			# TESTING: whether to include a group or not into a chart
			useGroup = switch xPattern
				when 'excluded' then false
				# else true
				# TODO: insert more patterns here ...
				when 'label', 'date', 'postCode', 'creditCard', 'currency' then true
				when 'intNumber' then switch xKeyPattern
					# TODO: insert key patterns here ... TODO
					when 'date' then true
					#when 'identifier', 'date' then true
					else false
				when 'floatNumber' then switch xKeyPattern
					when 'coefficient' then true
					else false
				else false
			# console.log useGroup
			#useGroup = true # TODO: Crap this when using patterns above ...

			if useGroup
				# TODO: mark here if there's already a count, it has no sense to do more than once ...
				noCount = true for group in $scope.groups when group.x is xKey
				if not noCount
					group = x: xKey, y: yKey, type: 'count', dimension: $scope.dimensions[i], group: null
					group.group = $scope.dimensions[i].group().reduceCount((d) -> d[yKey])
					$scope.groups.push group

				useSum = switch yPattern
					# TODO: discard more patterns here ....
					when 'label', 'date', 'postCode', 'creditCard', 'text' then false
					when 'intNumber' then switch yKeyPattern
						when 'identifier', 'date' then false
						else true
					else true
				#console.log useSum

				if useSum
					group2 = x: xKey, y: yKey, type: 'sum', dimension: $scope.dimensions[i], group: null
					group2.group = $scope.dimensions[i].group().reduceSum((d) -> d[yKey])
					$scope.groups.push group2

			return

		$scope.lineChart = (data, entry, xScale, ordinals) ->
			data["entry"] = entry
			data["xScale"] = xScale
			data["ordinals"] = ordinals

			data["colors"] = d3.scale.category20()
			data["colorAccessor"] = (d) -> parseInt(d.value)%20
			data["tickFormat"] = (d) =>
				switch $scope.data.patterns[entry.x].valuePattern
					when 'date' then d.getFullYear()
					when 'label', 'text' then "#{d}".substring(0, 20)
					else d

			data["xUnits"] = null
			data["tickValues"] = null
			if ordinals
				data["xUnits"] = dc.units.ordinal
				l = ordinals.length
				data["tickValues"] = [ordinals[0], ordinals[Math.floor(l/2)], ordinals[l-1]]

			data

		$scope.lineChartColorAccessor = (entry) ->
			parseInt(d.value) % 20

		$scope.lineChartXScale = (entry) ->
			switch $scope.data.patterns[entry.x].valuePattern
				# TODO: handle more patterns here .....
				when 'label'
					m = []
					m.push d.key for d in entry.group.all() when d not in m
					xScale = d3.scale.ordinal()
						.domain(m)
						.rangeBands([0, $scope.width])
				when 'date'
					xScale = d3.time.scale()
						#.domain(d3.extent(entry.group.all(), (d) -> PGPatternMatcher.parse(d.key, 'date')))
						.domain(d3.extent(entry.group.all(), (d) -> d.key))
						.range([0, $scope.width])
				else
					xScale = d3.scale.linear()
						.domain(d3.extent(entry.group.all(), (d) -> parseInt(d.key)))
						.range([0, $scope.width])

			xScale

		# $scope.drawLineChart = (entry, fixedId, xScale, ordinals) ->
		# 	console.log "drawLineChart", "##{fixedId}"
		# 	chart = dc.lineChart "##{fixedId}"

		# 	chart.width($scope.width)
		# 		.height($scope.height)
		# 		.margins({top: 10, right: 10, bottom: 30, left: 30})
		# 		.dimension(entry.dimension)
		# 		.group(entry.group)
		# 		.transitionDuration(500)
		# 		.colors(d3.scale.category20())
		# 		.colorAccessor((d) -> parseInt(d.value)%20)
		# 		.elasticY(true)
		# 		.x(xScale)
		# 		.xAxis()
		# 		.ticks(3)
		# 		.tickFormat(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.x].valuePattern
		# 					when 'date' then d.getFullYear()
		# 					when 'label', 'text' then "#{d}".substring(0, 20)
		# 					else d
		# 		)

		# 	if ordinals
		# 		chart.xUnits(dc.units.ordinal)
		# 		l = ordinals.length
		# 		chart.xAxis().tickValues([ordinals[0], ordinals[Math.floor(l/2)], ordinals[l-1]])
		# 		# TODO: Everything should deliver a chart, thrash the workaround below when well-tested
		# 		# if isNaN chart.yAxisMin()
		# 		#   $(@container).find("##{entry.type}-#{entry.x}-#{entry.y}").remove()
		# 		# else
		# 		#   chart.yAxis().ticks(3)

		# 	chart

		# $scope.drawBarsChart = (entry, fixedId, xScale, ordinals) ->
		# 	console.log "drawBarsChart"
		# 	chart = dc.barChart "##{fixedId}"
		# 	chart.width($scope.width)
		# 		.height($scope.height)
		# 		.margins({top: 10, right: 10, bottom: 30, left: 30})
		# 		.dimension(entry.dimension)
		# 		.group(entry.group)
		# 		.transitionDuration(500)
		# 		.centerBar(true)
		# 		.gap(2)
		# 		.colors(d3.scale.category20())
		# 		.colorAccessor((d) -> parseInt(d.data.value)%20)
		# 		.elasticY(true)
		# 		.x(xScale)
		# 		.xAxis()
		# 		.ticks(3)
		# 		.tickFormat(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.x].valuePattern
		# 					when 'date' then d.getFullYear()
		# 					when 'label', 'text' then "#{d}".substring(0, 20)
		# 					else d
		# 		)
		# 	if ordinals
		# 		chart.xUnits(dc.units.ordinal)
		# 		l = ordinals.length
		# 		chart.xAxis().tickValues([ordinals[0], ordinals[Math.floor(l/2)], ordinals[l-1]])
		# 	chart

		# $scope.drawRowsChart = (entry, fixedId) ->
		# 	console.log "drawRowsChart"
		# 	chart = dc.rowChart "##{fixedId}"
		# 	chart.width($scope.width)
		# 		.height($scope.height)
		# 		.margins({top: 5, right: 10, bottom: 20, left: 10})
		# 		.dimension(entry.dimension)
		# 		.group(entry.group)
		# 		.transitionDuration(500)
		# 		.gap(1)
		# 		.colors(d3.scale.category20())
		# 		.label(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.x].valuePattern
		# 					when 'date' then d.key.getFullYear()
		# 					when 'label', 'text' then "#{d.key}".substring(0, 20)
		# 					else d.key
		# 		)
		# 		.labelOffsetY($scope.height/(2*entry.group.size()))
		# 		.title((d) -> d.value)
		# 		.elasticX(true)
		# 		.xAxis()
		# 		.ticks(2)
		# 		.tickFormat(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.y].valuePattern
		# 					when 'date' then d.getFullYear()
		# 					when 'label', 'text' then "#{d}".substring(0, 20)
		# 					else d
		# 		)
		# 	chart

		# $scope.drawPieChart = (entry, fixedId) ->
		# 	console.log "drawPieChart"
		# 	chart = dc.pieChart "##{fixedId}"
		# 	chart.width($scope.width)
		# 		.height($scope.height)
		# 		.radius(Math.min($scope.width, $scope.height)/2)
		# 		.innerRadius(0.1*Math.min($scope.width, $scope.height))
		# 		.dimension(entry.dimension)
		# 		.group(entry.group)
		# 		.transitionDuration(500)
		# 		.colors(d3.scale.category20())
		# 		.label(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.x].valuePattern
		# 					when 'date' then d.data.key.getFullYear()
		# 					when 'label', 'text' then "#{d.data.key}".substring(0, 20)
		# 					else d.data.key
		# 		)
		# 		.minAngleForLabel(0.2)
		# 		.title((d) -> d.value)
		# 	chart

		# $scope.drawBubblesChart = (entry, fixedId, xScale) ->
		# 	console.log "drawBubblesChart"
		# 	svg = d3.select("##{fixedId}")
		# 		.append('svg')
		# 		.attr('width', $scope.width)
		# 		.attr('height', $scope.height)
		# 	chart = dc.bubbleOverlay("##{fixedId}")
		# 		.svg(svg)
		# 		.width($scope.width)
		# 		.height($scope.height)
		# 		.dimension(entry.dimension)
		# 		.group(entry.group)
		# 		.transitionDuration(500)
		# 		.keyAccessor((d) -> "Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_'))
		# 		.colors(d3.scale.category20())
		# 		.colorAccessor((d) -> parseInt(d.value)%20)
		# 		.radiusValueAccessor((d) -> d.value)
		# 		.r(d3.scale.linear().domain(d3.extent(entry.group.all(), (d) -> parseInt(d.value))))
		# 		.maxBubbleRelativeSize(0.1)
		# 		.label(
		# 			(d) =>
		# 				switch $scope.data.patterns[entry.x].valuePattern
		# 					when 'date' then d.key.getFullYear()
		# 					when 'label', 'text' then "#{d.key}".substring(0, 20)
		# 					else d.key
		# 		)
		# 		.minRadiusWithLabel(5)
		# 		.title((d) -> d.value)
		# 	chart.point(
		# 		"Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_')
		# 		0.1*$scope.width+0.8*xScale(d.key)
		# 		0.2*$scope.height+0.6*$scope.height*Math.random()
		# 	) for d in entry.group.all()
		# 	chart

		$scope.drawCharts = ->
			console.log "drawCharts"
			lastCharts = []

			$scope.chartsInfo = []

			for entry in $scope.groups
				do (entry) =>
					switch $scope.data.patterns[entry.x].valuePattern
						# TODO: handle more patterns here .....
						when 'label'
							m = []
							m.push d.key for d in entry.group.all() when d not in m
							xScale = d3.scale.ordinal()
								.domain(m)
								.rangeBands([0, $scope.width])
						when 'date'
							xScale = d3.time.scale()
								#.domain(d3.extent(entry.group.all(), (d) -> PGPatternMatcher.parse(d.key, 'date')))
								.domain(d3.extent(entry.group.all(), (d) -> d.key))
								.range([0, $scope.width])
						else
							xScale = d3.scale.linear()
								.domain(d3.extent(entry.group.all(), (d) -> parseInt(d.key)))
								.range([0, $scope.width])

					fixedId = "#{entry.type}-#{entry.x}-#{entry.y}".replace(/[^a-zA-Z0-9_-]/gi, '_')

					chartId = null
					for dcChart in $scope.charts
						do (dcChart) =>
							if not chartId and (
								not dcChart.maxEntries or entry.group.size() < dcChart.maxEntries
							) and dcChart.id not in lastCharts
								chartId = dcChart.id
								lastCharts.push dcChart.id
								return

					lastCharts = [] if lastCharts.length is $scope.charts.length

					plotChart = "line"
					if chartId in ['bar', 'pie', 'bubble']
						plotChart = chartId

					data =
						id: fixedId
						x: entry.x
						y: entry.y
						type: entry.type
						plot: plotChart + "Chart"

					chart = switch chartId
					# 	when 'row' then $scope.drawRowsChart entry, fixedId
					# 	when 'bar' then $scope.drawBarsChart entry, fixedId, xScale, m
					# 	when 'pie' then $scope.drawPieChart entry, fixedId
					# 	when 'bubble' then $scope.drawBubblesChart entry, fixedId, xScale
						when 'line' then data = $scope.lineChart data, entry, xScale, m
						else data = $scope.lineChart data, entry, xScale, m

					# chart = switch chartId
					# 	when 'row' then $scope.drawRowsChart entry, fixedId
					# 	when 'bar' then $scope.drawBarsChart entry, fixedId, xScale, m
					# 	when 'pie' then $scope.drawPieChart entry, fixedId
					# 	when 'bubble' then $scope.drawBubblesChart entry, fixedId, xScale
					# 	when 'line' then $scope.drawLineChart entry, fixedId, xScale, m
					# 	else $scope.drawLineChart entry, fixedId, xScale, m

					# # TESTING: How to get filtered data .... for maps or 3rd party elements
					# chart.on "filtered", (chart, filter) ->
					# 	console.log "filtered", chart.dimension().top(Infinity)
					# 	console.log filter
					# 	# Trigger 'update' for focusing maps on items bounds and 'updateOnlyItems' for no focus
					# 	# $(@).trigger 'updateOnlyItems', {elements: chart.dimension().bottom Infinity} #@todo
					# 	return

					$scope.chartsInfo.push data

					console.log "data:", fixedId, data

					return

			dc.renderAll()

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
