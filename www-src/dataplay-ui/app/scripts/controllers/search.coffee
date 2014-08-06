'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:SearchCtrl
 # @description
 # # SearchCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'SearchCtrl', ['$scope', '$location', '$routeParams', 'User', ($scope, $location, $routeParams, User) ->
		$scope.query = if $routeParams.query? then $routeParams.query else ""
		$scope.searchTimeout = null
		$scope.results = []

		$scope.search = () ->
			return if $scope.query.length < 3

			User.search($scope.query).success((data) ->
				$scope.results = data
				return
			).error (status, data) ->
				console.log "Search::search::Error:", status
				return

			return

		# Initiate search if we have /search/:query
		$scope.search()

		# Watch change to 'query' and do SILENT $location replace [No REFRESH]
		#	Note: Watch is initiated only after window.ready
		$scope.$watch "query", ((newVal, oldVal) ->
			if newVal.length >= 3
				qs = "/search/#{newVal}"
				$location.path qs, false
		), true

		return
	]
