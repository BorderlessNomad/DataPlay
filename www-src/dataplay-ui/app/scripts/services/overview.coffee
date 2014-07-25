'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Overview
 # @description
 # # Overview
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Overview', ['$http', 'Auth', 'config', ($http, Auth, config) ->
		reducedData: (guid, percent, min) ->
			$http.get config.api.base_url + "/getreduceddata" + "/" + guid + "/" + percent + "/" + min,
	]
