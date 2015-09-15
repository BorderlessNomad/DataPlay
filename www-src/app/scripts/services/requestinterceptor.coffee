'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.RequestInterceptor
 # @description
 # # RequestInterceptor
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.config ['$httpProvider', ($httpProvider) ->
		$httpProvider.interceptors.push "RequestInterceptor"
	]

	.factory 'RequestInterceptor', ['$q', 'Auth', 'config', ($q, Auth, config) ->
		"request": (reqConfig) ->
			reqConfig.headers[config.sessionHeader] = Auth.get config.sessionName

			reqConfig.url = config.api.base_url + "/" + reqConfig.url

			reqConfig || $q.when reqConfig
	]
