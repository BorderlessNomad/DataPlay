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

		$scope.search = (offset = 0, count = 9) ->
			return if $scope.query.length < 3
			User.search $scope.query, offset, count
				.success (data) ->
					if offset is 0
						$scope.results = data.Results
					else
						$scope.results = $scope.results.concat data.Results

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


		# Initiate search if we have /search/:query
		$scope.search()
		$scope.getNews()

		# Watch change to 'query' and do SILENT $location replace [No REFRESH]
		#	Note: Watch is initiated only after window.ready
		$scope.$watch "query", ((newVal, oldVal) ->
			if newVal.length >= 3
				qs = "/search/#{newVal}"
				$location.path qs, false
		), true

		return
	]
