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
			path = "chart/#{guid}/#{key}/#{type}/#{x}/#{y}"

			$http.get if z? then "#{path}/#{z}" else path

		correlated: (key) ->
			$http.get "chartcorrelated/#{key}"

		bookmark: (bookmarks) ->
			$http.post "setbookmark",
				data: bookmarks

		creditChart: (type, chartId, valFlag) ->
			path = 'chart/'
			path += if type is 'rid' then chartId.replace /\//g, '_' else chartId

			$http.put "#{path}/#{valFlag}"

		getObservations: (id) ->
			$http.get "observations/#{id}"

		createObservation: (did, x, y, message) ->
			$http.put "observations",
				did: '' + did
				x: "#{x}"
				y: "#{y}"
				comment: message

		creditObservation: (id, valFlag) ->
			$http.put "observations/#{id}/#{valFlag}"

		flagObservation: (id) ->
			$http.post "observations/flag/#{id}"

	]
