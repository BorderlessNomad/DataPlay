'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.ActivityMonitor
 # @description
 # # ActivityMonitor
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'ActivityMonitor', ['$http', 'Auth', 'config', ($http, Auth, config) ->
		get: (type = 'd') ->
			$http.get config.api.base_url + "/political/#{type}"
	]
