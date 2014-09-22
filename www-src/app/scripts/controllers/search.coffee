'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:SearchCtrl
 # @description
 # # SearchCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'SearchCtrl', ['$scope', '$location', '$routeParams', 'User', 'Overview', 'PatternMatcher', ($scope, $location, $routeParams, User, Overview, PatternMatcher) ->
		$scope.query = if $routeParams.query? then $routeParams.query else ""
		$scope.results = []
		$scope.rowedResults = [] # split into sub-arrays of 3

		$scope.totalResults = 0
		$scope.rowLimit = 3
		$scope.overview = []

		$scope.chartsRelated = []

		$scope.relatedChart = new RelatedCharts $scope.chartsRelated

		$scope.init = (reset = false) ->
			# Initiate search if we have /search/:query
			if reset
				$scope.chartsRelated = []
				$scope.totalResults
				$scope.relatedChart.chartsRelated = $scope.chartsRelated

			$scope.search()
			$scope.getNews()

		$scope.search = (offset = 0, count = 9) ->
			return if $scope.query.length < 3

			$scope.loading.related = true

			User.search $scope.query, offset, count
				.success (data) ->
					$scope.loading.related = false

					$scope.results = if offset is 0 then data.Results else $scope.results.concat data.Results

					$scope.results.forEach (r) ->
						r.graph = []
						r.error = null

						# Random
						# offset = Overview.getRandomInteger 0, 3
						$scope.getRelated r.GUID, 0

					$scope.rowedResults = $scope.splitIntoRows $scope.results
					$scope.totalResults = data.Total
					return
				.error (status, data) ->
					$scope.loading.related = false
					console.log "Search::search::Error:", status
					return
			return

		$scope.getNews = () ->
			return if $scope.query.length < 3

			User.getNews $scope.query
				.success (data) ->
					if data instanceof Array
						$scope.overview = data.map (item) ->
							date: Overview.humanDate new Date item.date
							title: item.title
							url: item.url
							thumbnail: item['image_url']

				.error (status, data) ->
					console.log "Search::getNews::Error:", status
					return

		$scope.splitIntoRows = (arr, numOfCols = 3) ->
			twoD = []
			for item, key in arr
				row = Math.floor key / numOfCols
				col = key % numOfCols
				if not twoD[row]
					twoD[row] = []
				twoD[row][col] = item
			return twoD

		$scope.showMore = ->
			$scope.rowLimit += 2
			if $scope.rowedResults.length < $scope.rowLimit
				# get more results
				$scope.search ($scope.rowLimit - 1) * 3, 9

		$scope.collapse = (item) ->
			item.show = false

		$scope.uncollapse = (item) ->
			item.show = true

		$scope.getRelated = (guid, offset) ->
			Overview.related guid, offset, 1
				.success (data) ->
					if data? and data.charts? and data.charts.length > 0
						for key, chart of data.charts
							continue unless $scope.relatedChart.isPlotAllowed chart.type

							key = parseInt(key)
							chart.guid = guid
							chart.key = key
							chart.id = "related-#{guid}-#{chart.key + $scope.offset.related}-#{chart.type}"
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

						console.log $scope.chartsRelated
					return
				.error (data, status) ->
					$scope.loading.related = false
					console.log "Overview::getRelated::Error:", status
					return

			return

		$scope.width = $scope.relatedChart.width
		$scope.height = $scope.relatedChart.height
		$scope.margin = $scope.relatedChart.margin
		$scope.loading = $scope.relatedChart.loading
		$scope.offset = $scope.relatedChart.offset
		$scope.limit = $scope.relatedChart.limit
		$scope.max = $scope.relatedChart.max

		$scope.hasRelatedCharts = $scope.relatedChart.hasRelatedCharts
		$scope.lineChartPostSetup = $scope.relatedChart.lineChartPostSetup
		$scope.rowChartPostSetup = $scope.relatedChart.rowChartPostSetup
		$scope.columnChartPostSetup = $scope.relatedChart.columnChartPostSetup
		$scope.pieChartPostSetup = $scope.relatedChart.pieChartPostSetup
		$scope.bubbleChartPostSetup = $scope.relatedChart.bubbleChartPostSetup

		return
	]
