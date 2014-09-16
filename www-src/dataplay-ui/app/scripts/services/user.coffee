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

		check: (username) ->
			$http.post config.api.base_url + "/user/check",
				username: username

		forgotPassword: (username) ->
			$http.post config.api.base_url + "/user/forgot",
				username: username

		token: (token, username, password) ->
			if password?
				$http.put config.api.base_url + "/user/reset/#{token}",
					username: username
					password: password
			else
				$http.get config.api.base_url + "/user/reset/#{token}/#{username}",

		resetPassword: (hash, password) ->
			$http.post config.api.base_url + "/user/reset",
				hash: hash
				password: password

		visited: () ->
			$http.get config.api.base_url + "/visited"

		search: (word, offset, count) ->
			path = "/search/#{word}"
			if offset?
				path += "/#{offset}"
				if count?
					path += "/#{count}"
			$http.get config.api.base_url + path
	]
