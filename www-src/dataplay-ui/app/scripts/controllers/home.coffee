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
		$scope.Auth = Auth
		$scope.username = Auth.get config.userName
		$scope.config = config
		$scope.lastVisited = []

		$scope.isActive = (path) ->
			$location.path().substr(0, path.length) is path

		$scope.getLastVisited = () ->
			unless Auth.isAuthenticated()
				$scope.lastVisited = []
				return

			User.visited().success((data) ->
				$scope.lastVisited = data
				return
			).error (status, data) ->
				console.log "Home::Visited::Error:", status
				return

			return

		return
	]
