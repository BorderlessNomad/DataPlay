'use strict'

###*
 # @ngdoc directive
 # @name dataplayApp.directive:match
 # @description
 # # match
###
angular.module("dataplayApp")
	.filter "orderObjectBy", [() ->
		(items, field, reverse) ->
			filtered = []
			angular.forEach items, (item) ->
				filtered.push item
				return

			filtered.sort (a, b) ->
				(if a[field] > b[field] then 1 else -1)

			filtered.reverse()  if reverse
			filtered
		]
