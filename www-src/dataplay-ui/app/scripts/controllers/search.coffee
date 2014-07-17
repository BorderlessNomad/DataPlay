'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:SearchCtrl
 # @description
 # # SearchCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'SearchCtrl', ['$scope', '$timeout', 'User', ($scope, $timeout, User) ->
		$scope.word = ""
		$scope.searchTimeout = null
		$scope.results = []

		$scope.search = () ->
			return if $scope.word.length < 3

			User.search($scope.word).success((data) ->
				$scope.results = data
				return
			).error (status, data) ->
				console.log "Search::search::Error:", status
				return

			return

		return
	]
