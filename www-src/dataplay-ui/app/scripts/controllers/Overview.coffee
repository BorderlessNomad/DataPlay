'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewCtrl
 # @description
 # # OverviewCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewCtrl', ['$scope', '$routeParams', 'Overview', ($scope, $routeParams, Overview) ->
		$scope.params = $routeParams
		$scope.info =
			name: $scope.params.id
			title: $scope.params.id

		$scope.info = () ->
			Overview.info $scope.params.id
				.success (data) ->
					$scope.info = data

					console.log $scope.info
					return
				.error (data, status) ->
					console.log "Overview::info::Error:", status
					return
	]
