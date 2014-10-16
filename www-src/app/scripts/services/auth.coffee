'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Auth
 # @description
 # # Auth
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Auth', ['ipCookie', 'config', (ipCookie, config) ->
		get: (key) ->
			token = ipCookie key
			if token? then token else false

		set: (key, value) ->
			ipCookie key, value,
				expires: 60 * 60 * 24 * 365, # 1 Year (seconds)
				expirationUnit: 'seconds'

		remove: (key) ->
			ipCookie.remove key

		isAuthenticated: () ->
			token = ipCookie config.sessionName
			if token? and token.length then true else false

		isAdmin: () ->
			true
	]

