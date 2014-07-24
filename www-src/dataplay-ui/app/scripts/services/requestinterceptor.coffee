'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.RequestInterceptor
 # @description
 # # RequestInterceptor
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.config ($httpProvider) ->
		$httpProvider.interceptors.push "RequestInterceptor"
		return

	.factory 'RequestInterceptor', ['$q', 'Auth', 'config', ($q, Auth, appConfig) ->
		"request": (config) ->
			# console.log "INTERCEPTED CONFIG", config, appConfig, Auth.isAuthenticated()
			config.headers[appConfig.sessionHeader] = Auth.isAuthenticated()

			config || $q.when config
	]
