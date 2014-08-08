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

		$scope.getRelatedCharts = () ->
			$scope.chartRegistryOffset = dc.chartRegistry.list().length

			Overview.related $scope.params.id
				.success (data) ->
					if data? and data.Charts? and data.Charts.length > 0
						for key, chart of data.Charts
							continue if chart.chart isnt 'line'
							chart.id = "#{$scope.params.id}-#{chart.xLabel}-#{chart.yLabel}-#{chart.chart}"

							chart.patterns = []
							chart.patterns[chart.xLabel] =
								valuePattern: PatternMatcher.getPattern chart.values[0]['x']
								keyPattern: PatternMatcher.getKeyPattern chart.values[0]['x']

							if chart.yLabel?
								chart.patterns[chart.yLabel] =
									valuePattern: PatternMatcher.getPattern chart.values[0]['y']
									keyPattern: PatternMatcher.getKeyPattern chart.values[0]['y']

							$scope.chartsInfo.push chart

						console.log "getRelatedCharts", $scope.chartsInfo

					return
				.error (data, status) ->
					console.log "Overview::getRelatedCharts::Error:", status
					return

			return

		$scope.getXScale = (data) ->
			xScale = switch data.patterns[data.xLabel].valuePattern
				# TODO: handle more patterns here .....
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.width]
				when 'date'
					d3.time.scale()
						.domain d3.extent(data.group.all(), (d) -> d.key)
						.range [0, $scope.width]
				else
					d3.scale.linear()
						.domain d3.extent(data.group.all(), (d) -> parseInt(d.key))
						.range [0, $scope.width]
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

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
