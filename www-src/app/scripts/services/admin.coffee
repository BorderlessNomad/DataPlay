'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Admin
 # @description
 # # Admin
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
	.factory 'Admin', ['$http', 'Auth', 'config', ($http, Auth, config) ->

		getUsers: (orderBy = 'uid', offset = 0, count = 15) ->
			order = ['uid', 'email', 'reputation', 'avatar', 'username', 'usertype', 'enabled']
			orderBy = 'uid' if orderBy not in order
			$http.get "admin/user/get/#{orderBy}/#{offset}/#{count}"

		editUser: (data) ->
			$http.put "admin/user/edit", data

		getObservations: (orderBy = 'observation_id', offset = 0, count = 15, flagged = false) ->
			order = ['comment', 'discovered_id', 'uid', 'rating', 'credited', 'discredited', 'observation_id', 'created', 'x', 'y', 'flagged']
			orderBy = 'observation_id' if orderBy not in order
			$http.get "admin/observations/get/#{orderBy}/#{offset}/#{count}/#{flagged}"

		deleteObservation: (id) ->
			$http.delete "admin/observations/#{id}"

	]
