'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCtrl
 # @description
 # # ChartsCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCtrl', ['$scope', '$routeParams', 'PatternMatcher', 'Charts', 'Tracker', ($scope, $routeParams, PatternMatcher, Charts, Tracker) ->
		$scope.params = $routeParams
		$scope.width = 1170;
		$scope.height = $scope.width * 9 / 16 # 16:9
		$scope.margin =
			top: 50
			right: 10
			bottom: 50
			left: 100
		$scope.chart =
			title: ""
			description: "N/A"
			data: null
		$scope.cfdata = null
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

		$scope.init = () ->
			# Track
			Tracker.visited $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y

			Charts.info $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y
				.success (data) ->
					if data?
						$scope.chart = data

						if data.desc? and data.desc.length > 0
							description = data.desc.replace /(h1>|h2>|h3>)/ig, 'h4>'
							description = description.replace /\n/ig, ''
							$scope.chart.description = description

						$scope.reduceData()

					console.log "Chart", $scope.chart

					return
				.error (data, status) ->
					console.log "Charts::getInfo::Error:", status
					return

		$scope.reduceData = () ->
			$scope.chart.data = $scope.chart.values
			$scope.chart.patterns = {}
			$scope.chart.patterns[$scope.chart.xLabel] =
				valuePattern: PatternMatcher.getPattern $scope.chart.values[0]['x']
				keyPattern: PatternMatcher.getKeyPattern $scope.chart.values[0]['x']

			if $scope.chart.patterns[$scope.chart.xLabel].valuePattern is 'date'
				for value, key in $scope.chart.values
					$scope.chart.values[key].x = new Date(value.x)

			if $scope.chart.yLabel?
				$scope.chart.patterns[$scope.chart.yLabel] =
					valuePattern: PatternMatcher.getPattern $scope.chart.values[0]['y']
					keyPattern: PatternMatcher.getKeyPattern $scope.chart.values[0]['y']

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
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

			chart.x $scope.getXScale data

			return

		$scope.rowChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> (i + 1) % 20

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
			data = $scope.chart

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
			data = $scope.chart

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

			chart.innerRadius 100

			chart.renderLabel false
			chart.label (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key} (#{Math.floor percent}%)"

			chart.renderTitle false
			chart.title (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key}: #{d.value} [#{Math.floor percent}%]"

			chart.legend dc.legend()

			chart.minAngleForLabel 0

			return

		$scope.bubbleChartPostSetup = (chart) ->
			# chart.colorAccessor (d, i) -> i + 1

			entry = crossfilter $scope.chart.data
			dimension = entry.dimension (d) -> d[0]
			group = dimension.group().reduceSum (d) -> d[1]

			chart.dimension dimension
			chart.group group

			svg = d3.select("#chart")
				.append 'svg'
				.attr 'width', $scope.width
				.attr 'height', $scope.height

			chart.svg svg
			chart.keyAccessor (d) -> "KEY_#{d.key}".replace(/\W+/g, "")
			chart.radiusValueAccessor (d) -> d.value
			chart.r d3.scale.linear().domain(d3.extent(group.all(), (d) -> parseInt(d.value)))
			chart.label (d) -> d.value
			chart.title (d) -> d.value

			ordinals = []
			ordinals.push d.key for d in group.all() when d not in ordinals
			xScale = $scope.getXScale ordinals

			for d in group.all()
				key = "KEY_#{d.key}".replace(/\W+/g, "")
				x = 0.1 * $scope.width + 0.8 * xScale d.key
				y = 0.2 * $scope.height + 0.6 * $scope.height * Math.random()

				chart.point key, x, y

			return

		return
	]
