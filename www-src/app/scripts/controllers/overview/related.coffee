'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewRelatedCtrl
 # @description
 # # OverviewRelatedCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewRelatedCtrl', ['$scope', '$routeParams', 'Overview', 'PatternMatcher', ($scope, $routeParams, Overview, PatternMatcher) ->
		$scope.allowed = ['line', 'bar', 'row', 'column', 'pie', 'bubble']
		$scope.params = $routeParams
		$scope.count = 3
		$scope.loading =
			related: false
		$scope.offset =
			related: 0
		$scope.limit =
			related: false
		$scope.max =
			related: 0
		$scope.chartsRelated = []

		$scope.xTicks = 6
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 0
			right: 0
			bottom: 0
			left: 0

		$scope.relatedChart = new RelatedCharts $scope.chartsRelated
		$scope.relatedChart.setPreview true
		$scope.relatedChart.setCustomMargin $scope.margin

		$scope.findById = (id) ->
			data = _.where($scope.chartsRelated,
				id: id
			)

			if data?[0]? then data[0] else null

		$scope.isPlotAllowed = (type) ->
			if type in $scope.allowed then true else false

		$scope.getRelatedCharts = () ->
			$scope.getRelated Overview.charts 'related'

			return

		$scope.hasRelatedCharts = () ->
			Object.keys($scope.chartsRelated).length

		$scope.getRelated = (count) ->
			$scope.loading.related = true

			if not count?
				count = $scope.max.related - $scope.offset.related
				count = if $scope.max.related and count < $scope.count then count else $scope.count

			Overview.related $scope.params.id, $scope.offset.related, count
				.success (data) ->
					$scope.loading.related = false

					if data? and data.charts? and data.charts.length > 0
						$scope.max.related = data.count

						for key, chart of data.charts
							continue unless $scope.relatedChart.isPlotAllowed chart.type

							chart.title = "#{chart.xLabel} vs #{chart.yLabel}"

							key = parseInt(key)
							chart.key = key
							chart.id = "related-#{$scope.params.id}-#{chart.key + $scope.relatedChart.offset.related}-#{chart.type}"
							chart.url = "charts/related/#{$scope.params.id}/#{chart.key}/#{chart.type}/#{chart.xLabel}/#{chart.yLabel}"
							chart.url += "/#{chart.zLabel}" if chart.type is 'bubble'

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

							$scope.chartsRelated.push chart

						console.log $scope.chartsRelated

						$scope.offset.related += count
						if $scope.offset.related >= $scope.max.related
							$scope.limit.related = true

						Overview.charts 'related', $scope.offset.related
					return
				.error (data, status) ->
					$scope.loading.related = false
					console.log "Overview::getRelated::Error:", status
					return

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()


		$scope.width = $scope.relatedChart.width
		$scope.height = $scope.relatedChart.height
		$scope.margin = $scope.relatedChart.margin

		$scope.hasRelatedCharts = $scope.relatedChart.hasRelatedCharts
		$scope.lineChartPostSetup = $scope.relatedChart.lineChartPostSetup
		$scope.rowChartPostSetup = $scope.relatedChart.rowChartPostSetup
		$scope.columnChartPostSetup = $scope.relatedChart.columnChartPostSetup
		$scope.pieChartPostSetup = $scope.relatedChart.pieChartPostSetup
		$scope.bubbleChartPostSetup = $scope.relatedChart.bubbleChartPostSetup

		return
	]
