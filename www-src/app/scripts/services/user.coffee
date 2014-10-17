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
			params =
				password: password
			if /^\S+@\S+\.\S+$/.test username
				params.email = username
			else
				params.username = username
			$http.post config.api.base_url + "/login", params

		logOut: (token) ->
			$http.delete config.api.base_url + "/logout"

		register: (username, email, password) ->
			$http.post config.api.base_url + "/register",
				username: username
				email: email
				password: password

		socialLogin: (data) ->
			$http.post config.api.base_url + "/login/social",
				'network': data['network'] or ''
				'id': data['id'] or ''
				'email': data['email'] or ''
				'full_name': data['full_name'] or ''
				'first_name': data['first_name'] or ''
				'last_name': data['last_name'] or ''
				'image': data['image'] or ''

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

		searchTweets: (word) ->
			$http.get config.api.base_url + "/tweets/#{word}"

		getNews: (query) ->
			if query instanceof Array
				query = query.join '_'
			$http.get config.api.base_url + "/news/search/#{query}"
	]
