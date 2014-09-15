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
		$scope.total = []

		$scope.overview = [
			{
				date: Overview.humanDate new Date
				title: "Health minister Jeremy Hunt announced new NHS budget for 2014"
				thumbnail: "http://i3.mirror.co.uk/incoming/article882107.ece/alternates/s615/Culture%20Secretary%20Jeremy%20Hunt%20looks%20on%20as%20Prime%20Minister%20David%20Cameron%20spoke%20during%20Prime%20Minister's%20Questions%20in%20the%20House%20of%20Commons"
			}
			{
				date: Overview.humanDate new Date
				title: "Ambulance response times improving year over year"
				thumbnail: "http://upload.wikimedia.org/wikipedia/commons/5/53/East_of_England_emergency_ambulance.jpg"
			}
		]

		$scope.search = () ->
			return if $scope.query.length < 3

			User.search($scope.query).success((data) ->
				$scope.results = data.Results
				$scope.total = data.Total

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
