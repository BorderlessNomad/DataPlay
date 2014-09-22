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

	.factory 'RequestInterceptor', ['$q', 'Auth', 'config', ($q, Auth, config) ->
		"request": (reqConfig) ->
			# console.log "INTERCEPTED CONFIG", reqConfig, config, Auth.isAuthenticated()
			reqConfig.headers[config.sessionHeader] = Auth.get config.sessionName

			reqConfig || $q.when reqConfig
	]
