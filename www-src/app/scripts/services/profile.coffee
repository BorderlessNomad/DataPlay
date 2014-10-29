'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.User
 # @description
 # # User
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Profile', ['$http', 'Auth', 'config', ($http, Auth, config) ->
		getInfo: (user) ->
			$http.get config.api.base_url + if user?.length > 0 then "/profile/#{user}" else "/user"

		setInfo: (email, username) ->
			$http.put config.api.base_url + "/user",
				email: email
				username: username

		getCreditDiscoveries: (user) ->
			$http.get config.api.base_url + if user?.length > 0 then "/profile/#{user}/credited" else "/profile/credited"

		getDiscoveries: (user) ->
			$http.get config.api.base_url + if user?.length > 0 then "/profile/#{user}/discoveries" else "/profile/discoveries"

		getObservations: (user) ->
			$http.get config.api.base_url + if user?.length > 0 then "/profile/#{user}/observations" else "/profile/observations"

	]
