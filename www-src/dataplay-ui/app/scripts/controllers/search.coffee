'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:SearchCtrl
 # @description
 # # SearchCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'SearchCtrl', ['$scope', '$location', '$routeParams', 'User', 'Overview', ($scope, $location, $routeParams, User, Overview) ->
		$scope.query = if $routeParams.query? then $routeParams.query else ""
		$scope.searchTimeout = null
		$scope.results = []
		$scope.rowedResults = [] # split into sub-arrays of 3

		$scope.totalResults = 0
		$scope.rowLimit = 3
		$scope.overview = []

		$scope.margin =
			top: 0
			right: 0
			bottom: 0
			left: 0

		$scope.init = () ->
			# Initiate search if we have /search/:query
			$scope.search()
			$scope.getNews()

		$scope.search = (offset = 0, count = 9) ->
			return if $scope.query.length < 3

			User.search $scope.query, offset, count
				.success (data) ->
					if offset is 0
						$scope.results = data.Results
					else
						$scope.results = $scope.results.concat data.Results

					$scope.results.forEach (r) ->
						r.graph = []
						r.error = null
						console.log r
						# Overview.related r.GUID
						# 	.success (graphdata) ->
						# 		console.log graphdata
						# 	.error (err, status) ->
						# 		r.error = switch
						# 			when err and err.message then err.message
						# 			when err and err.data and err.data.message then err.data.message
						# 			when err and err.data then err.data
						# 			when err then err
						# 			else ''
						# 		console.log "Search::search::getGraph::Error:", status

					$scope.rowedResults = $scope.splitIntoRows $scope.results
					$scope.totalResults = data.Total
					return
				.error (status, data) ->
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

		$scope.lineChartPostSetup = (details) ->
			(chart) ->
				data = details.graph

				return unless data and data.values and data.values.length > 0

				data.entry = crossfilter data.values
				data.dimension = data.entry.dimension (d) -> d.x
				data.group = data.dimension.group().reduceSum (d) -> d.y

				chart.dimension data.dimension
				chart.group data.group

				data.ordinals = []
				data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

				chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

				chart.xAxis().ticks 6

				chart.xAxisLabel false, 0
				chart.yAxisLabel false, 0

				chart.x $scope.getXScale data

				return

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

		# Watch change to 'query' and do SILENT $location replace [No REFRESH]
		#	Note: Watch is initiated only after window.ready
		$scope.$watch "query", ((newVal, oldVal) ->
			console.log newVal, newVal.length
			if newVal.length >= 3
				qs = "/search/#{newVal}"
				$location.path qs, false
		), true

		return
	]
