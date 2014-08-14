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
			Tracker.visited $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z

			Charts.info $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z
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
			chart.innerRadius 150

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.chart

			minR = null
			maxR = null

			###
			# TODO: when X/Y is String group all similar X/Y and SUM Y/X.
			# if Z is 'transaction_number' or 'invoice_number' ignore it
			# and replace radius with count of X/Y
			###

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
			data.ordinals.push d.key.split("|")[0] for d in data.group.all() when d not in data.ordinals
			console.log data.ordinals

			chart.keyAccessor (d) -> d.key.split("|")[0]
			chart.valueAccessor (d) -> d.key.split("|")[1]
			chart.radiusValueAccessor (d) ->
				r = Math.abs d.key.split("|")[2]
				if r >= minR then r else minR

			chart.x switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.width]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key.split("|")[0]
						.range [0, $scope.width]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[0]
						.range [0, $scope.width]

			chart.y switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.height]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key.split("|")[1]
						.range [0, $scope.height]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[1]
						.range [0, $scope.height]

			rScale = d3.scale.linear()
				.domain d3.extent data.group.all(), (d) -> Math.abs parseInt d.key.split("|")[2]
			chart.r rScale

			chart.label (d) -> x = d.key.split("|")[0]

			chart.title (d) ->
				x = d.key.split("|")[0]
				y = d.key.split("|")[1]
				z = d.key.split("|")[2]
				"#{data.xLabel}: #{x}\n#{data.yLabel}: #{y}\n#{data.zLabel}: #{z}"

			minRL = Math.log minR
			maxRL = Math.log maxR
			scale = Math.abs Math.log (maxRL - minRL) / (maxR - minR)

			console.log minR, maxR, (maxR - minR), minRL, maxRL, (maxRL - minRL), scale, scale / 100

			chart.maxBubbleRelativeSize scale / 100

			chart.zoomOutRestrict true
			chart.mouseZoomable true

			# chart.renderlet (c) ->
			# 	circles = c.svg().selectAll('g.chart-body').selectAll('g circle')
			# 	for circle in circles[0]
			# 		r = circle.r
			# 		if r.baseVal.value < 0
			# 			r.baseVal.value = Math.abs r.baseVal.value
			# 			r.animVal = r.baseVal

			return

		return
	]
