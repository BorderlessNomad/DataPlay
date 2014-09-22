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
		getInfo: ->
			$http.get config.api.base_url + "/user"

		setInfo: (email, username) ->
			$http.put config.api.base_url + "/user",
				email: email
				username: username

		getValidDiscoveries: ->
			$http.get config.api.base_url + "/profile/validated"

		getDiscoveries: ->
			$http.get config.api.base_url + "/profile/discoveries"

		getObservations: ->
			$http.get config.api.base_url + "/profile/observations"

	]
