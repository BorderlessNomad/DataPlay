'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Overview
 # @description
 # # Overview
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Overview', ['$http', 'config', ($http, config) ->
		chart = {}
		chartRegistryOffset = 0
		monthNames = [
			"Jan"
			"Feb"
			"Mar"
			"Apr"
			"May"
			"Jun"
			"Jul"
			"Aug"
			"Sep"
			"Oct"
			"Nov"
			"Dec"
		]
		reducedData: (guid, percent, min) ->
			$http.get config.api.base_url + "/getreduceddata/#{guid}/#{percent}/#{min}"
		related: (guid, offset, count) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			$http.get config.api.base_url + "/related/#{guid}/#{offset}/#{count}"
		correlated: (guid, offset, count) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			$http.get config.api.base_url + "/correlated/#{guid}/0/0"
		charts: (key, value) ->
			unless key?
				chart = {}
				return null
			chart[key] = value if value?
			return chart[key] if chart[key]?
			null
		getChartOffset: (chart) ->
			flag = if chart.__dc_flag__? then chart.__dc_flag__ else 0
			console.log "getChartOffset", flag, chartRegistryOffset, dc.chartRegistry.list().length, flag - chartRegistryOffset - 1
			flag - chartRegistryOffset - 1
		updateChartRegistry: (offset) ->
			chartRegistryOffset = offset

	]
