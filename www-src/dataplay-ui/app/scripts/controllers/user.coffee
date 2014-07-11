'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:UserCtrl
 # @description
 # # UserCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'UserCtrl', ['$scope', '$location', 'ipCookie', 'User', 'Auth', 'config', ($scope, $location, ipCookie, User, Auth, config) ->
		$scope.user =
			username: null
			password: null
			password_confirm: null
			message: null

		$scope.login = (user) ->
			if user.username? and user.password?
				User.logIn(user.username, user.password).success((data) ->
					$scope.processLogin data

					return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::Login::Error:", status
					return

			return

		$scope.processLogin = (data) ->
			Auth.username = data.username

			ipCookie data.session.name, data.session.value,
				expires: data.session.expiry
				expirationUnit: 'seconds'

			$location.path "/home"

			return

		$scope.logout = () ->
			token = Auth.isAuthenticated()

			if token isnt false
				User.logOut(token).success((data) ->
					Auth.username = null

					ipCookie.remove config.cookieName

					$location.path "/"
					return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::Logout::Error:", status
					return

			$location.path "/login"

			return

		$scope.register = (user) ->
			if user.username? and user.password?
				User.register(user.username, user.password).success((data) ->
					$scope.processLogin data

					return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::Register::Error:", status
					return

			return

		return
	]
