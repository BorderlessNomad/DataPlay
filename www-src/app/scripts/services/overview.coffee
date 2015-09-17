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
		chart = {}
		monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]

		charts: (key, value) ->
			unless key?
				chart = {}
				return null
			chart[key] = value if value?
			return chart[key] if chart[key]?
			null

		humanDate: (date) ->
			"#{date.getDate()} #{monthNames[date.getMonth()]}, #{date.getFullYear()}"

		getRandomInteger: (min, max) ->
			Math.floor(Math.random() * (max - min) + min)

		info: (guid) ->
			$http.get "chartinfo/#{guid}"

		related: (guid, offset, count) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			$http.get "related/#{guid}/#{offset}/#{count}"

		correlated: (guid, offset, count, depth) ->
			offset = if offset? then offset else 0
			count = if count? then count else 3
			depth = if depth? then depth else 100
			$http.get "correlated/#{guid}/true/#{offset}/#{count}"
	]
