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

		$scope.getRelatedCharts = () ->
			$scope.chartRegistryOffset = dc.chartRegistry.list().length

			allowed = ['line', 'bar', 'row', 'column', 'pie', 'bubble']

			Overview.related $scope.params.id
				.success (data) ->
					if data? and data.Charts? and data.Charts.length > 0
						for key, chart of data.Charts
							if chart.type not in allowed
								continue

							chart.id = "#{$scope.params.id}-#{chart.xLabel}-#{chart.yLabel}-#{chart.type}"

							chart.patterns = []
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
						.domain d3.extent data.group.all(), (d) -> parseInt(d.key)
						.range [0, $scope.width]

			xScale

		$scope.getYScale = (data) ->
			xScale = switch data.patterns[data.yLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.height]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key
						.range [0, $scope.height]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt(d.key)
						.range [0, $scope.height]

			xScale

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
			data.dimension = data.entry.dimension (d) -> d.y
			data.group = data.dimension.group().reduceSum (d) -> d.x

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.xAxis()
				.ticks $scope.xTicks

			chart.x $scope.getXScale data

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
			data.dimension = data.entry.dimension (d) -> d.x
			data.groupSum = 0
			data.group = data.dimension.group().reduceSum (d) ->
				data.groupSum += parseFloat(d.y)
				d.y

			chart.dimension data.dimension
			chart.group data.group

			chart.colorAccessor (d, i) -> i + 1

			# chart.innerRadius 100

			chart.renderLabel false
			chart.label (d) -> "#{d.key} (#{Math.floor d.value / data.groupSum * 100}%)"

			chart.renderTitle false
			chart.title (d) -> "#{d.key}: #{d.value} [#{Math.floor d.value / data.groupSum * 100}%]"

			# chart.legend dc.legend()

			# chart.minAngleForLabel 0

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.chartsInfo[$scope.getChartOffset chart]

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			svg = d3.select("#chart")
				.append 'svg'
				.attr 'width', $scope.width
				.attr 'height', $scope.height

			chart.svg svg
			chart.keyAccessor (d) -> "KEY_#{d.key}".replace(/\W+/g, "")
			chart.radiusValueAccessor (d) -> d.value
			chart.r d3.scale.linear().domain(d3.extent(data.group.all(), (d) -> parseInt(d.value)))
			chart.label (d) -> d.value
			chart.title (d) -> d.value

			ordinals = []
			ordinals.push d.key for d in data.group.all() when d not in ordinals
			xScale = $scope.getXScale ordinals

			for d in data.group.all()
				key = "KEY_#{d.key}".replace(/\W+/g, "")
				x = 0.1 * $scope.width + 0.8 * xScale d.key
				y = 0.2 * $scope.height + 0.6 * $scope.height * Math.random()

				chart.point key, x, y

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
