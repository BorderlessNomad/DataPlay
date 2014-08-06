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
		info: (guid) ->
			$http.get config.api.base_url + "/getinfo/#{guid}"

		reducedData: (guid, x, y, percent, min) ->
			$http.get config.api.base_url + "/getreduceddata/#{guid}/#{x}/#{y}/#{percent}/#{min}"

		groupedData: (guid, x, y) ->
			$http.get config.api.base_url + "/getdatagrouped/#{guid}/#{x}/#{y}"

		identifyData: (guid) ->
			$http.get config.api.base_url + "/identifydata/#{guid}"

		bookmark: (bookmarks) ->
			$http.post config.api.base_url + "/setbookmark",
				data: bookmarks
	]
