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

		getUsers: (orderby = 'uid', offset = 0, count = 15) ->
			orderbys = ['uid', 'email', 'reputation', 'avatar', 'username', 'usertype', 'enabled']
			if orderbys.indexOf(orderby) is -1 then orderby = 'uid'
			$http.get config.api.base_url + "/admin/user/get/#{orderby}/#{offset}/#{count}"

		editUser: (data) ->
			$http.put config.api.base_url + "/admin/user/edit", data

	]
