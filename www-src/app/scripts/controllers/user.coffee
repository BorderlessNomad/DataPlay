'use strict'

hello.init {
	facebook: '849402518433862'
	# windows : '000000004812EF76'
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
			type: null
			login:
				token: $scope.params.token
				username: null
				password: null
				message: null
			register:
				token: $scope.params.token
				username: null
				email: null
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
					$scope.user.type = 'login'
					$scope.user.login.message = status
					console.log "User::Login::Error:", status
					return

			return

		$scope.processLogin = (data) ->
			Auth.set config.userId, data.uid
			Auth.set config.userName, data.user
			Auth.set config.userType, data.usertype
			Auth.set config.sessionName, data.session

			$location.path "/home"

			return

		$scope.socialLogin = (type) ->
			hello(type).login()
				.then (auth) ->
					hello(auth.network).api('/me')
						.then (data) ->
							core =
								network: type
								id: data.id or ''
								email: data.email or ''
								"full_name": data.name or ''
								"first_name": data['first_name'] or ''
								"last_name": data['last_name'] or ''

							if type is 'facebook'
								core.image = data.picture or ''
							else if type is 'google'
								core.image = data.image?.url or ''

							User.socialLogin core
								.success (res) ->
									$scope.processLogin res
								.error (msg) ->
									$scope.user.type = 'socialLogin'
									$scope.user.login.message = msg
									console.log "User::SocialLogin::Error:", msg
									return

						, (e) -> console.log "Error on /me", e
				, (e) -> console.log "Error on login", e


		$scope.logout = () ->
			token = Auth.get config.sessionName

			if token isnt false
				User.logOut(token).success((data) ->
					Auth.remove config.userId
					Auth.remove config.userName
					Auth.remove config.userType
					Auth.remove config.sessionName

					$location.path "/"
					return
				).error (status, data) ->
					$scope.user.type = 'logout'
					$scope.user.login.message = status
					console.log "User::Logout::Error:", status
					return

			$location.path "/user/login"

			return

		$scope.register = (user) ->
			if user.username? and user.email? and user.password?
				User.register(user.username, user.email, user.password).success((data) ->
					$scope.processLogin data

					return
				).error (status, data) ->
					$scope.user.type = 'register'
					$scope.user.register.message = status
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
						$scope.user.login.token = data.token
						$scope.user.register.token = data.token
						return
					).error (status, data) ->
						$scope.forgotPassword.valid = false
						$scope.user.type = 'forgotPassword'
						$scope.user.login.message = status
						console.log "User::ForgotPassword::token::Error:", status, data
						return
				).error (status, data) ->
					$scope.user.type = 'forgotPassword'
					$scope.user.login.message = status
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
					$scope.user.type = 'resetPasswordCheck'
					$scope.user.login.message = status
					console.log "User::ResetPassword::check::Error:", status, data
					return

		$scope.resetPassword = (user) ->
			$scope.closeAlert()

			if user.token? and user.username? and user.password?
				User.token(user.token, user.username, user.password).success((data) ->
					$scope.resetPassword.saved = true

					return
				).error (status, data) ->
					$scope.user.type = 'resetPassword'
					$scope.user.login.message = status
					console.log "User::ResetPassword::save::Error:", status, data
					return

		$scope.hasError = (type) ->
			$scope.user[type]?.message?.length && $scope.user.type is type

		$scope.closeAlert = () ->
			$scope.user.type = null
			$scope.user.login.message = null
			$scope.user.register.message = null

		$scope.changeTab = (tab) ->
			$scope.currentTab = tab

		return
	]
