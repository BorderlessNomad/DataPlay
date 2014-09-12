'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewCorrelatedCtrl
 # @description
 # # OverviewCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewCorrelatedCtrl', ['$scope', '$routeParams', 'Overview', 'PatternMatcher', ($scope, $routeParams, Overview, PatternMatcher) ->
		$scope.allowed = ['line', 'bar', 'row', 'column', 'bubble', 'scatter', 'stacked']
		$scope.params = $routeParams
		$scope.count = 3
		$scope.loading =
			correlated: false
		$scope.offset =
			correlated: 0
		$scope.limit =
			correlated: false
		$scope.max =
			correlated: 0
		$scope.chartsCorrelated = {}

		$scope.xTicks = 6
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 70

		$scope.isPlotAllowed = (type) ->
			if type in $scope.allowed then true else false

		$scope.getCorrelatedCharts = () ->
			$scope.getCorrelated Overview.charts 'correlated'

			return

		$scope.hasCorrelatedCharts = () ->
			Object.keys($scope.chartsCorrelated).length

		$scope.getCorrelated = (count) ->
			$scope.loading.correlated = true

			if not count?
				count = $scope.max.correlated - $scope.offset.correlated
				count = if $scope.max.correlated and count < $scope.count then count else $scope.count

			depth = if $scope.max.correlated then 0 else 100

			Overview.correlated $scope.params.id, $scope.offset.correlated, count, depth
				.success (data) ->
					if data? and data.charts? and data.charts.length > 0
						$scope.loading.correlated = false

						$scope.max.correlated = data.count

						for key, chart of data.charts
							continue unless $scope.isPlotAllowed chart.type
							# continue unless chart.type is 'line'

							key = parseInt(key)
							chart.key = key
							chart.id = "correlated-#{$scope.params.id}-#{chart.key + $scope.offset.correlated}-#{chart.type}"
							chart.url = "/charts/correlated/#{$scope.params.id}/#{chart.correlationid}/#{chart.type}/#{chart.table1.xLabel}/#{chart.table1.yLabel}"
							chart.url += "/#{chart.table1.zLabel}" if chart.type is 'bubble' or chart.type is 'scatter'

							chart.patterns = {}
							chart.patterns[chart.table1.xLabel] =
								valuePattern: PatternMatcher.getPattern chart.table1.values[0]['x']
								keyPattern: PatternMatcher.getKeyPattern chart.table1.values[0]['x']

							if chart.patterns[chart.table1.xLabel].valuePattern is 'date'
								for value, key in chart.table1.values
									chart.table1.values[key].x = new Date(value.x)

								for value, key in chart.table2.values
									chart.table2.values[key].x = new Date(value.x)

							if chart.table1.yLabel?
								chart.patterns[chart.table1.yLabel] =
									valuePattern: PatternMatcher.getPattern chart.table1.values[0]['y']
									keyPattern: PatternMatcher.getKeyPattern chart.table1.values[0]['y']

							$scope.chartsCorrelated[chart.id] = chart if PatternMatcher.includePattern(
								chart.patterns[chart.table1.xLabel].valuePattern,
								chart.patterns[chart.table1.xLabel].keyPattern
							)

						console.log $scope.chartsCorrelated

						$scope.offset.correlated += count
						if $scope.offset.correlated >= $scope.max.correlated
							$scope.limit.correlated = true

						Overview.charts 'correlated', $scope.offset.correlated
					return
				.error (data, status) ->
					$scope.loading.correlated = false
					console.log "Overview::getCorrelated::Error:", status
					return

			return

		$scope.getXScale = (data) ->
			xScale = switch data.patterns[data.table1.xLabel].valuePattern
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

		$scope.getXUnits = (data) ->
			xUnits = switch data.patterns[data.table1.xLabel].valuePattern
				when 'date' then d3.time.years
				when 'intNumber' then dc.units.integers
				when 'label', 'text' then dc.units.ordinal
				else dc.units.ordinal

			xUnits

		$scope.getYScale = (data) ->
			yScale = switch data.patterns[data.table1.xLabel].valuePattern
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
			data = $scope.chartsCorrelated[chart.anchorName()]

			data.xDomain = [new Date(data.from), new Date(data.to)]
			data.entry = crossfilter data.table1.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			data.entry2 = crossfilter data.table2.values
			data.dimension2 = data.entry2.dimension (d) -> d.x
			data.group2 = data.dimension2.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group
			# chart.stack data.group2

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

			chart.xAxis().ticks $scope.xTicks

			xScale = d3.time.scale()
				.domain data.xDomain
				.range [0, $scope.width]
			chart.x xScale

			return

		$scope.rowChartPostSetup = (chart) ->
			data = $scope.chartsCorrelated[chart.anchorName()]

			data.entry = crossfilter data.table1.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			data.entry2 = crossfilter data.table2.values
			data.dimension2 = data.entry2.dimension (d) -> d.x
			data.group2 = data.dimension2.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.xAxis().ticks $scope.xTicks

			chart.x $scope.getYScale data

			# chart.xUnits $scope.getXUnits data if data.ordinals?.length > 0

			return

		$scope.columnChartPostSetup = (chart) ->
			data = $scope.chartsCorrelated[chart.anchorName()]

			data.entry = crossfilter data.table1.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			data.entry2 = crossfilter data.table2.values
			data.dimension2 = data.entry2.dimension (d) -> d.x
			data.group2 = data.dimension2.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group
			# chart.stack data.group2

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.x $scope.getXScale data

			chart.xUnits $scope.getXUnits data if data.ordinals?.length > 0

			return

		$scope.stackedChartPostSetup = (chart) ->
			data = $scope.chartsCorrelated[chart.anchorName()]

			data.entry = crossfilter data.table1.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			data.entry2 = crossfilter data.table2.values
			data.dimension2 = data.entry2.dimension (d) -> d.x
			data.group2 = data.dimension2.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group
			# chart.stack data.group2

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.x $scope.getXScale data

			chart.xUnits $scope.getXUnits data if data.ordinals?.length > 0

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.chartsCorrelated[chart.anchorName()]

			minR = null
			maxR = null

			data.entry = crossfilter data.table1.values
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

			chart.x switch data.patterns[data.table1.xLabel].valuePattern
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

			chart.y switch data.patterns[data.table1.xLabel].valuePattern
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

			# chart.label (d) -> x = d.key.split("|")[0]

			chart.title (d) ->
				x = d.key.split("|")[0]
				y = d.key.split("|")[1]
				z = d.key.split("|")[2]
				"#{data.table1.xLabel}: #{x}\n#{data.table1.yLabel}: #{y}\n#{data.table1.zLabel}: #{z}"

			minRL = Math.log minR
			maxRL = Math.log maxR
			scale = Math.abs Math.log (maxRL - minRL) / (maxR - minR)

			chart.maxBubbleRelativeSize scale / 100

			return

		$scope.scatterChartPostSetup = (chart) ->
			data = $scope.chartsCorrelated[chart.anchorName()]

			data.entry = crossfilter data.table1.values
			data.dimension = data.entry.dimension (d) -> "#{d.x}|#{d.y}|#{d.z}"
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			for d in data.group.all() when d not in data.ordinals
				data.ordinals.push d.key.split("|")[0]

			chart.keyAccessor (d) -> d.key.split("|")[0]
			chart.valueAccessor (d) -> d.key.split("|")[1]

			chart.x switch data.patterns[data.table1.xLabel].valuePattern
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

			chart.y switch data.patterns[data.table1.xLabel].valuePattern
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

			chart.label (d) -> x = d.key.split("|")[0]

			chart.title (d) ->
				x = d.key.split("|")[0]
				y = d.key.split("|")[1]
				"#{data.table1.xLabel}: #{x}\n#{data.table1.yLabel}: #{y}"

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
