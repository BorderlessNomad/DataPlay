'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:AdminUsersCtrl
 # @description
 # # AdminUsersCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'AdminUsersCtrl', ['$scope', '$location', 'Auth', 'Overview', 'config', ($scope, $location, Auth, Overview, config) ->

		$scope.modal =
			shown: false
			content:
				title: null
				type: null
				item: null

		$scope.headers = [
			{key: 'id', display: '#'}
			{key: 'avatar', display: 'Avatar'}
			{key: 'username', display: 'Username'}
			{key: 'reputation', display: 'Reputation'}
			{key: 'isAdmin', display: 'Admin?'}
		]

		$scope.users = [
			{
				id: 180
				username: 'jack'
				avatar: 'https://pbs.twimg.com/profile_images/3164870237/efe0014851567f9dca856297f8292bf1.jpeg'
				reputation: 25
				usertype: 0
			}
			{
				id: 11
				username: 'mayur'
				avatar: 'http://www.gravatar.com/avatar/848f09d47991c7995cb7ba9bbf3e8b93?d=identicon'
				reputation: 207
				usertype: 1
			}
			{
				id: 181
				username: 'glyn'
				avatar: 'http://www.gravatar.com/avatar/9f1839175aab93c0a0fd9e36623fe17d?d=identicon'
				reputation: 9001
				usertype: 0
			}
		]

		$scope.init = () ->
			if not $scope.isAdmin()
				$location.path '/home'

		$scope.isAdmin = () ->
			true # TODO: actually check whether current user is an admin

		$scope.showModal = (title, type, item) ->
			if item?
				$scope.modal.content.title = title
				$scope.modal.content.type = type
				$scope.modal.content.item = item

				$scope.modal.shown = true
				$('#admin-modal').modal 'show'

		$scope.closeModal = () ->
			$scope.modal.shown = false
			$('#admin-modal').modal 'hide'
			return



		$scope.view = (user) ->
			console.log '   '
			$scope.showModal 'Viewing', 'view', user
			return

		$scope.edit = (user) ->
			console.log '   '
			$scope.showModal 'Editing', 'edit', user
			return

		$scope.disable = (user) ->
			console.log '   '
			$scope.showModal 'Disabling', 'disable', user
			return

		return
	]
