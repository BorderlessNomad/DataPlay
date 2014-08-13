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
		$scope.allowed = ['line', 'bar', 'row', 'column', 'pie', 'bubble']
		# $scope.allowed = ['bubble']
		$scope.params = $routeParams
		$scope.chartsInfo = []
		$scope.chartRegistryOffset = 0

		$scope.xTicks = 6
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 70
		$scope.monthNames = [
			"Jan"
			"Feb"
			"Mar"
			"Apr"
			"May"
			"Jun"
			"Jul"
			"Aug"
			"Sep"
			"Oct"
			"Nov"
			"Dec"
		]

		$scope.getChartOffset = (chart) ->
			chart.__dc_flag__ - $scope.chartRegistryOffset - 1

		$scope.isPlotAllowed = (type) ->
			if type in $scope.allowed then true else false

		$scope.getRelatedCharts = () ->
			$scope.chartRegistryOffset = dc.chartRegistry.list().length

			Overview.related $scope.params.id
				.success (data) ->
					if data? and data.Charts? and data.Charts.length > 0
						for key, chart of data.Charts
							continue unless $scope.isPlotAllowed chart.type

							chart.id = "#{$scope.params.id}-#{chart.xLabel}-#{chart.yLabel}-#{chart.type}"

							chart.patterns = {}
							chart.patterns[chart.xLabel] =
								valuePattern: PatternMatcher.getPattern chart.values[0]['x']
								keyPattern: PatternMatcher.getKeyPattern chart.values[0]['x']

							if chart.patterns[chart.xLabel].valuePattern is 'date'
								for value, key in chart.values
									chart.values[key].x = new Date(value.x)

							if chart.yLabel?
								chart.patterns[chart.yLabel] =
									valuePattern: PatternMatcher.getPattern chart.values[0]['y']
									keyPattern: PatternMatcher.getKeyPattern chart.values[0]['y']

							if $scope.includePattern chart
								$scope.chartsInfo.push chart

						console.log "getRelatedCharts", $scope.chartsInfo, $scope.chartsInfo.length

					return
				.error (data, status) ->
					console.log "Overview::getRelatedCharts::Error:", status
					return

			return

		$scope.includePattern = (data) ->
			xPattern = data.patterns[data.xLabel].valuePattern
			xKeyPattern = data.patterns[data.xLabel].keyPattern

			useGroup = switch xPattern
				when 'excluded' then false
				when 'label', 'date', 'postCode', 'creditCard', 'currency' then true
				when 'intNumber', 'floatNumber' then switch xKeyPattern
					when 'date', 'coefficient' then true
					when null then true
					else false
				else false

			useGroup

		$scope.getXScale = (data) ->
			xScale = switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.width]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key
						.range [0, $scope.width]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key
						.range [0, $scope.width]

			xScale

		$scope.getYScale = (data) ->
			yScale = switch data.patterns[data.yLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.height]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.value
						.range [0, $scope.height]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.value
						.range [0, $scope.height]
						.nice()

			yScale

		$scope.lineChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

			chart.xAxis()
				.ticks $scope.xTicks

			chart.x $scope.getXScale data

			return

		$scope.rowChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.xAxis()
				.ticks $scope.xTicks

			chart.x $scope.getYScale data

			if ordinals? and ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			return

		$scope.columnChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.x $scope.getXScale data

			if ordinals? and ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			return

		$scope.pieChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) ->
				if data.patterns[data.xLabel].valuePattern is 'date'
					return "#{d.x.getDate()} #{$scope.monthNames[d.x.getMonth()]} #{d.x.getFullYear()}"
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
			chart.label (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key} (#{Math.floor percent}%)"

			chart.renderTitle false
			chart.title (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key}: #{d.value} [#{Math.floor percent}%]"

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			tempData = []

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) ->
				tempData["#{d.x}_#{d.y}"] = d.z
				d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.keyAccessor (d) -> d.key
			chart.valueAccessor (d) -> d.value
			chart.radiusValueAccessor (d) -> tempData["#{d.key}_#{d.value}"]

			chart.x $scope.getXScale data

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
