'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:UserCtrl
 # @description
 # # UserCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'UserCtrl', ($scope, $location, ngCookies, User, Auth) ->

		$scope.user =
			username: null
			password: null

		$scope.login = (user) ->
			if user.username? and user.password?
				User.logIn(user.username, user.password).success((data) ->
					Auth.isLogged = true
					# store Cookie
					$location.path "/home"
					return
				).error (status, data) ->
					console.log "User::Login::Error:", status, data
					return

			return

		$scope.logout = () ->
			if Auth.isLogged
				# Send DELETE request to Server
				User.logOut(token).success((data) ->
					Auth.isLogged = false
					# delete Cookie
					$location.path "/"
					return
				).error (status, data) ->
					console.log "User::Logout::Error:", status, data
					return

			return

		$scope.register = (user) ->
			if user.username? and user.password?
				User.register(user.username, user.password).success((data) ->
					$location.path "/login"
					return
				).error (status, data) ->
					console.log "User::Register::Error:", status, data
					return

			return

		return
