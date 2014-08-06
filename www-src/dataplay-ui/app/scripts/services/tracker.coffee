'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Tracker
 # @description
 # # Tracker
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Tracker', ['$http', 'config', ($http, config) ->
		visited: (guid, type, x, y) ->
			$http.post config.api.base_url + "/visited",
				guid: guid
				info:
					type: type
					x: x
					y: y
	]
