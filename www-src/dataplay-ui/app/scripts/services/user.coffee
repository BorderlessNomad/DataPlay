'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.User
 # @description
 # # User
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Auth', (ipCookie, config) -> {
		username: null

		isAuthenticated: () ->
			token = ipCookie config.cookieName
			if token?
				token
			else
				false
	}

	.factory 'User', ($http, config) -> {
		logIn: (username, password) ->
			$http.post config.api.base_url + "/login",
				username: username
				password: password

		logOut: (token) ->
			$http.delete config.api.base_url + "/logout" + "/" + token

		register: (username, password) ->
			$http.post config.api.base_url + "/register",
				username: username
				password: password
	}
