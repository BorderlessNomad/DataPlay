'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:HomeCtrl
 # @description
 # # HomeCtrl
 # Controller of the dataplayApp
###
d3.selection::duration = -> @
# d3.selection::transition = -> @

angular.module('dataplayApp')
	.controller 'HomeCtrl', ['$scope', '$location', 'Home', 'Auth', 'Overview', 'PatternMatcher', 'config', ($scope, $location, Home, Auth, Overview, PatternMatcher, config) ->
		$scope.config = config

		$scope.searchquery = ''

		$scope.creditPatterns = null

		$scope.myActivity = null
		$scope.recentObservations = null
		$scope.dataExperts = null

		$scope.loading =
			charts: true

		$scope.chartsRelated = []

		$scope.relatedChart = new RelatedCharts $scope.chartsRelated
		$scope.relatedChart.setPreview true

		$scope.init = ->
			$scope.loading.charts = true;
			Home.getAwaitingCredit()
				.success (data) ->
					$scope.loading.charts = false;
					if data? and data.charts? and data.charts.length > 0
						for key, chart of data.charts
							continue unless $scope.relatedChart.isPlotAllowed chart.type
							continue unless key < 4

							key = parseInt(key)

							if chart.relationid?
								guid = chart.relationid.split("/")[0]

								chart.key = key
								chart.id = "related-#{guid}-#{chart.key + $scope.relatedChart.offset.related}-#{chart.type}"
								chart.url = "charts/related/#{guid}/#{chart.key}/#{chart.type}/#{chart.xLabel}/#{chart.yLabel}"
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

								$scope.chartsRelated.push chart if PatternMatcher.includePattern(
									chart.patterns[chart.xLabel].valuePattern,
									chart.patterns[chart.xLabel].keyPattern
								)

							else if chart.correlationid?
								chartObj = new CorrelatedChart chart.type

								if not chartObj.error
									chartObj.info =
										key: key
										id: "correlated-#{chart.correlationid}"
										url: "charts/correlated/#{chart['source_title']}/#{chart.correlationid}/#{chart.type}/#{chart['source_X']}/#{chart['source_Y']}"
										title: [chart.table1.title, chart.table2.title]
									chartObj.info.url += "/#{chart.table1.zLabel}" if chart.type is 'bubble'

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

									$scope.chartsRelated.push chartObj

						console.log $scope.chartsRelated
					else
						$scope.creditPatterns = []
				.error ->
					$scope.creditPatterns = []

			Home.getActivityStream()
				.success (data) ->
					if data instanceof Array
						$scope.myActivity = data.map (d) ->
							date: Overview.humanDate new Date d.time
							pretext: d.activitystring
							linktext: d.patternid
							url: d.linkstring
					else
						$scope.myActivity = []
				.error ->
					$scope.myActivity = []

			Home.getRecentObservations()
				.success (data) ->
					if data instanceof Array
						$scope.recentObservations = data.map (d) ->
							user:
								name: d.username
								avatar: d.avatar or "http://www.gravatar.com/avatar/#{d.MD5email}?d=identicon"
							text: d.comment
							url: d.linkstring
					else
						$scope.recentObservations = []
				.error ->
					$scope.recentObservations = []

			Home.getDataExperts()
				.success (data) ->
					if data instanceof Array

						medals = ['gold', 'silver', 'bronze']

						$scope.dataExperts = data.map (d, key) ->
							obj =
								rank: key + 1
								name: d.username
								avatar: d.avatar or "http://www.gravatar.com/avatar/#{d.MD5email}?d=identicon"
								score: d.reputation

							if obj.rank <= 3 then obj.rankclass = medals[obj.rank - 1]

							obj
					else
						$scope.dataExperts = []
				.error ->
					$scope.dataExperts = []

		$scope.search = ->
			$location.path "/search/#{$scope.searchquery}"

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
