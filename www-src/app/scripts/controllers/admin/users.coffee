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
			actions:
				view: 'Viewing'
				edit: 'Editing'
				disable: 'Disabling'

		$scope.headers = [
			{key: 'id', display: '#'}
			{key: 'avatar', display: 'Avatar'}
			{key: 'username', display: 'Username'}
			{key: 'reputation', display: 'Reputation'}
			{key: 'isAdmin', display: 'Admin?'}
		]

		$scope.pagination =
			perPage: 15
			pageNumber: 1

		# $scope.users = [
		# 	{ id: 1, username: 'jack', avatar: 'https://pbs.twimg.com/profile_images/3164870237/efe0014851567f9dca856297f8292bf1.jpeg', reputation: 25, usertype: 0 }
		# 	{ id: 2, username: 'mayur', avatar: 'http://www.gravatar.com/avatar/848f09d47991c7995cb7ba9bbf3e8b93?d=identicon', reputation: 207, usertype: 1}
		# 	{ id: 3, username: 'glyn', avatar: 'http://www.gravatar.com/avatar/9f1839175aab93c0a0fd9e36623fe17d?d=identicon', reputation: 9001, usertype: 0 }
		# ]

		$scope.users = do ->
			total = 27
			result = []
			samples = [
				{ id: 1, username: 'jack', avatar: 'https://pbs.twimg.com/profile_images/3164870237/efe0014851567f9dca856297f8292bf1.jpeg', reputation: 25, usertype: 0 }
				{ id: 2, username: 'mayur', avatar: 'http://www.gravatar.com/avatar/848f09d47991c7995cb7ba9bbf3e8b93?d=identicon', reputation: 207, usertype: 1}
				{ id: 3, username: 'glyn', avatar: 'http://www.gravatar.com/avatar/9f1839175aab93c0a0fd9e36623fe17d?d=identicon', reputation: 9001, usertype: 0 }
			]
			for item in [1..total]
				newItem = _.clone samples[Math.floor(Math.random() * samples.length)]
				newItem.id = item
				newItem.reputation += Math.floor(Math.random() * 30)
				newItem.usertype = Math.floor(Math.random() * 2)
				result.push newItem
			result


		# General controls
		$scope.init = () ->
			if not $scope.isAdmin()
				$location.path '/home'

		$scope.isAdmin = () ->
			true # TODO: actually check whether current user is an admin

		$scope.showModal = (type, item) ->
			if item?
				$scope.modal.content.type = type
				$scope.modal.content.item = _.clone item

				$scope.modal.content.item.usertype = !! item.usertype

				$scope.modal.shown = true
				$('#admin-modal').modal 'show'
			return

		$scope.closeModal = () ->
			$scope.modal.shown = false
			$('#admin-modal').modal 'hide'
			return

		$scope.submitForm = () ->
			if $scope.modal.content.type is 'edit'
				$scope.edit()
			else if $scope.modal.content.type is 'disable'
				$scope.disable()

		$scope.edit = () ->
			console.log "Save the edit"
			return

		$scope.disable = () ->
			console.log "Disable the user"
			return



		# Pagination
		$scope.totalPages = (total) ->
			Math.ceil total / $scope.pagination.perPage

		$scope.range = (usrsLen) ->
			end = $scope.totalPages usrsLen
			[1..end]

		$scope.changePage = (page, add) ->
			newVal = if add? then $scope.pagination.pageNumber + page else page
			if newVal > 0 and newVal <= $scope.totalPages $scope.users.length
				$scope.pagination.pageNumber = newVal

		return
	]
