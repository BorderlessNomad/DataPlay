'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Home
 # @description
 # # Home
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Home', ['$http', 'Auth', 'config', ($http, Auth, config) ->

		getStats: () ->
			$http.get "home/data"

		getTopRated: () ->
			$http.get "chart/toprated"

		getAwaitingCredit: () ->
			$http.get "chart/awaitingcredit"

		getActivityStream: () ->
			$http.get "user/activitystream"

		getRecentObservations: () ->
			$http.get "recentobservations"

		getDataExperts: () ->
			$http.get "user/experts"
	]
