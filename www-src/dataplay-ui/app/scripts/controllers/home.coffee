'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:HomeCtrl
 # @description
 # # HomeCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'HomeCtrl', ['$scope', '$location', 'User', 'Auth', 'config', ($scope, $location, User, Auth, config) ->
		$scope.user = Auth
		$scope.config = config

		$scope.isActive = (path) ->
			$location.path().substr(0, path.length) is path

		return
	]
