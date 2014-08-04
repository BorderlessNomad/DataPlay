'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.User
 # @description
 # # User
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'User', ['$http', 'Auth', 'config', ($http, Auth, config) ->
		logIn: (username, password) ->
			$http.post config.api.base_url + "/login",
				username: username
				password: password

		logOut: (token) ->
			$http.delete config.api.base_url + "/logout"

		register: (username, password) ->
			$http.post config.api.base_url + "/register",
				username: username
				password: password

		visited: () ->
			$http.get config.api.base_url + "/visited"

		search: (word) ->
			$http.get config.api.base_url + "/search/#{word}"
	]
