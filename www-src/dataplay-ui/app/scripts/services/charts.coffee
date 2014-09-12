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
		related: (guid, key, type, x, y, z) ->
			path = "/chart/#{guid}/#{key}/#{type}/#{x}/#{y}"
			if z? then path += "/#{z}"
			$http.get config.api.base_url + path

		correlated: (key) ->
			$http.get config.api.base_url + "/chartcorrelated/#{key}"

		bookmark: (bookmarks) ->
			$http.post config.api.base_url + "/setbookmark",
				data: bookmarks

		validateChart: (chartId, valFlag) ->
			path = "/chart/#{chartId}"

			if valFlag
				path += "/#{valFlag}"

			$http.put config.api.base_url + path

		getObservations: (id) ->
			$http.get config.api.base_url + "/observations/#{id}"

		createObservation: (did, x, y, message) ->
			$http.put config.api.base_url + "/observations",
				did: did
				x: "#{x}"
				y: "#{y}"
				comment: message

		validateObservation: (id, valFlag) ->
			$http.put config.api.base_url + "/observations/#{id}/#{valFlag}"
	]
