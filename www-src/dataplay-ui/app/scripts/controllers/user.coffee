'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:UserCtrl
 # @description
 # # UserCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'UserCtrl', ['$scope', '$location', 'User', 'Auth', 'config', ($scope, $location, User, Auth, config) ->
		$scope.user =
			username: null
			password: null
			password_confirm: null
			message: null

		$scope.forgotPassword =
			valid: false
			sent: false

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
			Auth.set config.userName, data.user
			Auth.set config.sessionName, data.session

			$location.path "/home"

			return

		$scope.logout = () ->
			token = Auth.get config.sessionName

			if token isnt false
				User.logOut(token).success((data) ->
					Auth.remove config.sessionName
					Auth.remove config.userName

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

		$scope.forgotPassword = (user) ->
			$scope.closeAlert()

			if user.username?
				# User.check user.username
				# 	.then (result) ->
				# 		$scope.forgotPassword.valid = true
				# 		User.forgotPassword user.username
				# 	.then (result) ->
				# 		$scope.forgotPassword.sent = true
				# 	.error (result) ->
				# 		console.log "error", result

				User.check(user.username).success((data) ->
					$scope.forgotPassword.valid = true

					User.forgotPassword(user.username).success((data) ->
						$scope.forgotPassword.sent = true

						return
					).error (status, data) ->
						$scope.forgotPassword.valid = false
						$scope.user.message = status
						console.log "User::ForgotPassword::token::Error:", status, data
						return
				).error (status, data) ->
					$scope.forgotPassword.valid = false
					$scope.user.message = status
					console.log "User::ForgotPassword::check::Error:", status, data
					return

			return

		$scope.hasError = () ->
			if $scope.user.message?.length then true else false

		$scope.closeAlert = () ->
			$scope.user.message = null

		return
	]
