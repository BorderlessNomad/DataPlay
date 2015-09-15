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
			$http.post "login", params

		logOut: (token) ->
			$http.delete "logout"

		register: (username, email, password) ->
			$http.post "register",
				username: username
				email: email
				password: password

		socialLogin: (data) ->
			$http.post "login/social",
				'network': data['network'] or ''
				'id': data['id'] or ''
				'email': data['email'] or ''
				'full_name': data['full_name'] or ''
				'first_name': data['first_name'] or ''
				'last_name': data['last_name'] or ''
				'image': data['image'] or ''

		check: (email) ->
			$http.post "user/check",
				email: email

		forgotPassword: (email) ->
			$http.post "user/forgot",
				email: email

		token: (token, email, password) ->
			if password?
				$http.put "user/reset/#{token}",
					email: email
					password: password
			else
				$http.get "user/reset/#{token}/#{email}"

		resetPassword: (hash, password) ->
			$http.post "user/reset",
				hash: hash
				password: password

		visited: () ->
			$http.get "visited"

		search: (word, offset, count) ->
			word = word.replace(/\/|\\/g, ' ')
			path = "/search/#{word}"
			if offset?
				path += "/#{offset}"
				if count?
					path += "/#{count}"
			$http.get path

		searchTweets: (word) ->
			$http.get "tweets/#{word}"

		getNews: (query) ->
			if query instanceof Array
				query = query.join '_'
			query = query.replace(/\s{1,}|\%20|\/|\\/g, '_')
			$http.get "news/search/#{query}"
	]
