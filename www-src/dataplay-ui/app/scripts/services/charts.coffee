'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Charts
 # @description
 # # Charts
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Charts', ['$http', 'config', ($http, config) ->
		info: (guid, key, type, x, y, z) ->
			if z?
				$http.get config.api.base_url + "/chart/#{guid}/#{key}/#{type}/#{x}/#{y}/#{z}"
			else
				$http.get config.api.base_url + "/chart/#{guid}/#{key}/#{type}/#{x}/#{y}"

		bookmark: (bookmarks) ->
			$http.post config.api.base_url + "/setbookmark",
				data: bookmarks
	]
