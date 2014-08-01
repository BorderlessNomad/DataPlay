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
		$scope.chartRegistryOffset = 0

		$scope.guid = null
		$scope.data = null
		$scope.keys = []
		$scope.cfdata = null
		$scope.dimensions = []
		$scope.groups = []
		$scope.charts = [
			{ id: 'row', maxEntries: 15 }
			{ id: 'bar', maxEntries: 30 }
			{ id: 'pie', maxEntries: 45 }
			{ id: 'bubble', maxEntries: 60 }
			{ id: 'line' }
		]
		$scope.xTicks = 6
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 70

		$scope.getChartOffset = (chart) ->
			chart.__dc_flag__ - $scope.chartRegistryOffset - 1

		$scope.getCharts = () ->
			$scope.chartRegistryOffset = dc.chartRegistry.list().length

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

						# console.log "Patterns: ", patterns

						$scope.plot $scope.params.id, {dataset: data, patterns: patterns}

					return
				.error (data) ->
					console.log "Overview::getCharts::Error:", status
					return

			return

		$scope.plot = (guid, data, width, height) ->
			$scope.guid = guid
			$scope.data = data

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
				noCount = true for group in $scope.groups when group.x is xKey
				if not noCount
					data =
						x: xKey
						y: yKey
						type: 'count'
						dimension: $scope.dimensions[i]
						group: null

					data.group = $scope.dimensions[i].group().reduceCount((d) -> d[yKey])

					$scope.groups.push data

				useSum = switch yPattern
					# TODO: discard more patterns here ....
					when 'label', 'date', 'postCode', 'creditCard', 'text' then false
					when 'intNumber' then switch yKeyPattern
						when 'identifier', 'date' then false
						else true
					else true

				if useSum
					data =
						x: xKey
						y: yKey
						type: 'sum'
						dimension: $scope.dimensions[i]
						group: null

					data.group = $scope.dimensions[i].group().reduceSum((d) -> d[yKey])

					$scope.groups.push data

			return

		$scope.drawLineChart = (data, entry) ->
			# Line chart calculations (pre-render) goes here

			data

		$scope.lineChartPostSetup = (chart) ->
			# Line chart manipulations (post-render) goes here
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			chart.colorAccessor (d) -> parseInt(d.value) % 20

			chart.xAxis()
				.ticks $scope.xTicks
				# .tickFormat (d) =>
				# 	switch $scope.data.patterns[data.entry.x].valuePattern
				# 		when 'date' then new Date(d)
				# 		when 'intNumber' then parseInt(d)
				# 		when 'label', 'text' then d.substring 0, 20
				# 		else d

			if data.ordinals? and data.ordinals.length > 0
				chart.xUnits switch $scope.data.patterns[data.entry.x].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

				# l = data.ordinals.length
				# tickValues = []
				# tickValues.push data.ordinals[0]
				# tickValues.push data.ordinals[Math.floor(l / $scope.xTicks * i)] for i in [1..$scope.xTicks - 1]
				# tickValues.push data.ordinals[l - 1]
				# chart.xAxis().tickValues tickValues

			return

		$scope.drawBarChart = (data, entry) ->

			data

		$scope.barChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			chart.colorAccessor (d) -> parseInt(d.value) % 20

			chart.xAxis()
				.ticks $scope.xTicks
				# .tickFormat (d) =>
				# 	switch $scope.data.patterns[data.entry.x].valuePattern
				# 		when 'date' then d.getFullYear()
				# 		when 'label', 'text' then d.substring 0, 20
				# 		else d

			if data.ordinals? and data.ordinals.length > 0
				chart.xUnits switch $scope.data.patterns[data.entry.x].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal
				l = data.ordinals.length
				# chart.xAxis().tickValues [data.ordinals[0], data.ordinals[Math.floor(l/2)], data.ordinals[l-1]]

			return

		$scope.drawRowChart = (data, entry) ->
			data["labelOffsetY"] = $scope.height / ( 2 * entry.group.size())
			data["colorAccessor"] = (d) -> parseInt(d.value) % 20
			data["tickFormat"] = (d) => # needs xAxis()
				switch $scope.data.patterns[entry.x].valuePattern
					when 'date' then d.getFullYear()
					when 'label', 'text' then d.substring 0, 20
					else d
			data["label"] = (d) => # needs xAxis()
				switch $scope.data.patterns[entry.x].valuePattern
					when 'date' then d.data.key.getFullYear()
					when 'label', 'text' then d.data.key.substring 0, 20
					else d.data.key
			data["title"] = (d) -> d.value

			data

		$scope.drawPieChart = (data, entry) ->
			data["radius"] = Math.min($scope.width, $scope.height) / 2
			data["innerRadius"] = 0.1 * Math.min $scope.width, $scope.height
			data["colorAccessor"] = (d) -> parseInt(d.value) % 20
			data["label"] = (d) =>
				switch $scope.data.patterns[entry.x].valuePattern
					when 'date' then d.data.key.getFullYear()
					when 'label', 'text' then d.data.key.substring 0, 20
					else d.data.key
			data["title"] = (d) -> d.value

			data

		$scope.drawBubbleChart = (data, entry) ->

			data

		$scope.bubbleChartPostSetup = (chart, options) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			svg = d3.select("##{data.id}")
				.append 'svg'
				.attr 'width', $scope.width
				.attr 'height', $scope.height

			chart
				.svg svg
				.keyAccessor (d) -> "Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_')
				.colorAccessor (d) -> parseInt(d.value) % 20
				.radiusValueAccessor (d) -> d.value
				.r d3.scale.linear().domain(d3.extent(data.entry.group.all(), (d) -> parseInt(d.value)))
				.label (d) =>
					switch $scope.data.patterns[data.entry.x].valuePattern
						when 'date' then d.key.getFullYear()
						when 'label', 'text' then d.key.substring 0, 20
						else d.key
				.title (d) -> d.value

			chart.point(
				"Key#{d.key}".replace(/[^a-zA-Z0-9_-]/gi, '_')
				0.1 * $scope.width + 0.8 * data.xScale d.key
				0.2 * $scope.height + 0.6 * $scope.height * Math.random()
			) for d in data.entry.group.all()

			return

		$scope.transformChartType = (chart, x, y) ->
			if chart is 'line'
				if not $scope.data.patterns[x].keyPattern?
					return 'bar'

			if chart is 'bar'
				if $scope.data.patterns[x].keyPattern? and $scope.data.patterns[x].keyPattern is 'date'
					return 'line'

			chart

		$scope.drawCharts = ->
			lastCharts = []

			$scope.chartsInfo = []

			for entry in $scope.groups
				do (entry) =>
					ordinals = []
					ordinals.push d.key for d in entry.group.all() when d not in ordinals

					xScale = switch $scope.data.patterns[entry.x].valuePattern
						# TODO: handle more patterns here .....
						when 'label'
							d3.scale.ordinal()
								.domain ordinals
								.rangeBands [0, $scope.width]
						when 'date'
							d3.time.scale()
								.domain d3.extent(entry.group.all(), (d) -> d.key)
								.range [0, $scope.width]
						else
							d3.scale.linear()
								.domain d3.extent(entry.group.all(), (d) -> parseInt(d.key))
								.range [0, $scope.width]

					fixedId = "#{entry.type}-#{entry.x}-#{entry.y}".replace(/[^a-zA-Z0-9_-]/gi, '_')

					chartId = null
					for chart in $scope.charts
						do (chart) =>
							if not chartId and (not chart.maxEntries or entry.group.size() < chart.maxEntries) and chart.id not in lastCharts
								chartId = chart.id
								lastCharts.push chart.id
								return

					lastCharts = [] if lastCharts.length is $scope.charts.length

					plotChart = 'line'
					if chartId in ['bar', 'pie', 'bubble']
						plotChart = chartId

					plotChart = $scope.transformChartType plotChart, entry.x, entry.y

					xUnits = dc.units.integers
					if ordinals? and ordinals.length
						xUnits = switch $scope.data.patterns[entry.x].valuePattern
							when 'date' then d3.time.years
							when 'intNumber' then dc.units.integers
							when 'label', 'text' then dc.units.ordinal
							else dc.units.integers

					data =
						id: fixedId
						x: entry.x
						y: entry.y
						type: entry.type
						plot: plotChart
						xScale: xScale
						ordinals: ordinals
						xUnits: xUnits
						entry: entry

					data = switch chartId
						when 'line'
							$scope.drawLineChart data, entry
						when 'bar'
							$scope.drawBarChart data, entry
						when 'row'
							$scope.drawRowChart data, entry
						when 'pie'
							$scope.drawPieChart data, entry
						when 'bubble'
							$scope.drawBubbleChart data, entry
						else
							$scope.drawLineChart data, entry

					$scope.chartsInfo.push data

					return

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
