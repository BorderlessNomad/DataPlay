'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:RecentCtrl
 # @description
 # # RecentCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'RecentCtrl', ['$scope', '$location', 'User', 'Auth', 'Overview', 'config', ($scope, $location, User, Auth, Overview, config) ->
		$scope.Auth = Auth
		$scope.username = Auth.get config.userName
		$scope.config = config
		$scope.lastVisited = []

		$scope.isActive = (path) ->
			$location.path().substr(0, path.length) is path

		$scope.getLastVisited = () ->
			# Reset Overview cache
			Overview.charts null

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
