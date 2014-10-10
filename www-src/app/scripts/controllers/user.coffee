'use strict'

hello.init {
	facebook: '849402518433862'
	twitter: 'd1rdbGQ4Hd6ZvURsaxAKcgUiW'
	google  : '233668141851-tjicormqo6m2ld7rdckk2ok2vnqf3ff3.apps.googleusercontent.com'
}, {
	redirect_uri: 'redirect.html'
	scope: 'email'
}

###*
 # @ngdoc function
 # @name dataplayApp.controller:UserCtrl
 # @description
 # # UserCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'UserCtrl', ['$scope', '$location', '$routeParams', 'User', 'Auth', 'config', ($scope, $location, $routeParams, User, Auth, config) ->
		$scope.params = $routeParams
		$scope.user =
			token: $scope.params.token
			username: null
			password: null
			password_confirm: null
			message: null

		$scope.forgotPassword =
			valid: false
			sent: false

		$scope.resetPassword =
			valid: false
			saved: false

		$scope.currentTab = 'details'

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

		$scope.socialLogin = (type) ->
			hello(type).login()
				.then (auth) ->
					hello(auth.network).api('/me')
						.then (data) ->
							if type is 'facebook'
								obj =
									id: data.id or ''
									email: data.email or ''
									name:
										full: data.name or ''
										first: data['first_name'] or ''
										last: data['last_name'] or ''
									image: data.picture or ''
							else if type is 'google'
								obj =
									id: data.id or ''
									email: data.email or ''
									name:
										full: data.name or data.displayName or ''
										first: data['first_name'] or ''
										last: data['last_name'] or ''
									image: data.image?.url or ''
						, (e) -> console.log "Error on /me", e
				, (e) -> console.log "Error on login", e


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

			$location.path "/user/login"

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
				User.check(user.username).success((data) ->
					$scope.forgotPassword.valid = true

					User.forgotPassword(user.username).success((data) ->
						$scope.forgotPassword.sent = true
						$scope.user.token = data.token
						return
					).error (status, data) ->
						$scope.forgotPassword.valid = false
						$scope.user.message = status
						console.log "User::ForgotPassword::token::Error:", status, data
						return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::ForgotPassword::check::Error:", status, data
					return

			return

		$scope.resetPasswordCheck = (user) ->
			$scope.closeAlert()

			if user.token? and user.username?
				User.token(user.token, user.username).success((data) ->
					$scope.resetPassword.valid = true

					return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::ResetPassword::check::Error:", status, data
					return

		$scope.resetPassword = (user) ->
			$scope.closeAlert()

			if user.token? and user.username? and user.password?
				User.token(user.token, user.username, user.password).success((data) ->
					$scope.resetPassword.saved = true

					return
				).error (status, data) ->
					$scope.user.message = status
					console.log "User::ResetPassword::save::Error:", status, data
					return

		$scope.hasError = () ->
			if $scope.user.message?.length then true else false

		$scope.closeAlert = () ->
			$scope.user.message = null

		$scope.changeTab = (tab) ->
			$scope.currentTab = tab

		return
	]
