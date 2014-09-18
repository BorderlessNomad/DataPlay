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
		timeFormatter = d3.time.format.multi([
			[".%L", (d) -> d.getMilliseconds()]
			[":%S", (d) -> d.getSeconds()]
			["%I:%M", (d) -> d.getMinutes()]
			["%I %p", (d) -> d.getHours()]
			["%a %d", (d) -> d.getDay() and d.getDate() isnt 1]
			["%b %d", (d) -> d.getDate() isnt 1]
			["%B", (d) -> d.getMonth()]
			["%Y", (d) -> true]
		])
		reducedData: (guid, percent, min) ->
			$http.get config.api.base_url + "/getreduceddata/#{guid}/#{percent}/#{min}"
		related: (guid, offset, count) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			$http.get config.api.base_url + "/related/#{guid}/#{offset}/#{count}"
		correlated: (guid, offset, count, depth) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			depth = if depth? then depth else 100
			$http.get config.api.base_url + "/correlated/#{guid}/#{offset}/#{count}/#{depth}"
		charts: (key, value) ->
			unless key?
				chart = {}
				return null
			chart[key] = value if value?
			return chart[key] if chart[key]?
			null
		humanDate: (date) ->
			"#{date.getDate()} #{monthNames[date.getMonth()]}, #{date.getFullYear()}"
	]
