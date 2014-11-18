'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewCorrelatedCtrl
 # @description
 # # OverviewCorrelatedCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewCorrelatedCtrl', ['$scope', '$routeParams', 'Overview', 'PatternMatcher', ($scope, $routeParams, Overview, PatternMatcher) ->
		$scope.allowed = ['line', 'bar', 'row', 'column', 'bubble', 'scatter', 'stacked']
		$scope.params = $routeParams
		$scope.count = 6
		$scope.loading =
			correlated: false
		$scope.offset =
			correlated: 0
		$scope.limit =
			correlated: false
		$scope.max =
			correlated: 0
		$scope.chartsCorrelated = []

		$scope.xTicks = 6
		$scope.width = 350
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 70

		$scope.findById = (id) ->
			data = _.where($scope.chartsCorrelated,
				id: id
			)

			if data?[0]? then data[0] else null

		$scope.isPlotAllowed = (type) ->
			if type in $scope.allowed then true else false

		$scope.getCorrelatedCharts = () ->
			$scope.getCorrelated Overview.charts 'correlated'
			return

		$scope.hasCorrelatedCharts = () ->
			Object.keys($scope.chartsCorrelated).length

		$scope.getCorrelated = (count) ->
			$scope.loading.correlated = true

			if not count? or count is 0
				count = $scope.max.correlated - $scope.offset.correlated
				count = if $scope.max.correlated and count < $scope.count then count else $scope.count
			else
				generate = true

			generate = if generate or $scope.max.correlated then false else true

			Overview.correlated $scope.params.id, $scope.offset.correlated, count, false #generate
				.success (data) ->
					$scope.loading.correlated = false

					if data? and data.charts? and data.charts.length > 0
						$scope.max.correlated = data.count

						for key, chart of data.charts
							continue unless $scope.isPlotAllowed chart.type
							# continue unless chart.type is 'line'

							chartObj = new CorrelatedChart chart.type

							chartObj.title = chart.table2.title
							chartObj.coeff = Math.floor Math.abs chart.coefficient * 100

							if not chartObj.error
								key = parseInt(key)
								chartObj.info =
									key: key
									id: "correlated-#{$scope.params.id}-#{chart.key + $scope.offset.correlated}-#{chart.type}"
									url: "charts/correlated/#{chart.correlationid}"

								[1..2].forEach (i) ->
									vals = chartObj.translateData chart['table' + i].values, chart.type
									dataRange = do ->
										min = d3.min vals, (item) -> parseFloat item.y
										[
											if min > 0 then 0 else min
											d3.max vals, (item) -> parseFloat item.y
										]
									type = if chart.type is 'column' or chart.type is 'bar' then 'bar' else 'area'

									chartObj.data.push
										key: chart['table' + i].title
										type: type
										yAxis: i
										values: vals
									chartObj.options.chart['yDomain' + i] = dataRange
									chartObj.options.chart['yAxis' + i].tickValues = [0]
									chartObj.options.chart.xAxis.tickValues = []

								chartObj.setAxisTypes 'none', 'none', 'none'
								chartObj.setSize null, 200
								chartObj.setMargin 25, 25, 25, 25
								chartObj.setLegend false
								chartObj.setTooltips false
								chartObj.setPreview true
								chartObj.setLabels chart

								$scope.chartsCorrelated.push chartObj

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


		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
