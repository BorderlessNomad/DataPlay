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
			$http.get if user?.length > 0 then "profile/#{user}" else "user"

		setInfo: (email, username) ->
			$http.put "user",
				email: email
				username: username

		getCreditDiscoveries: (user) ->
			$http.get if user?.length > 0 then "profile/#{user}/credited" else "profile/credited"

		getDiscoveries: (user) ->
			$http.get if user?.length > 0 then "profile/#{user}/discoveries" else "profile/discoveries"

		getObservations: (user) ->
			$http.get if user?.length > 0 then "profile/#{user}/observations" else "profile/observations"

	]
