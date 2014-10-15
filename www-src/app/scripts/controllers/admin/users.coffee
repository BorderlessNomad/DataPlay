'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:AdminUsersCtrl
 # @description
 # # AdminUsersCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'AdminUsersCtrl', ['$scope', '$location', 'Admin', 'Overview', 'config', ($scope, $location, Admin, Overview, config) ->

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

		$scope.pagination =
			perPage: 10
			pageNumber: 1
			total: 0

		$scope.users = []


		# General controls
		$scope.init = () ->
			if not $scope.isAdmin()
				$location.path '/home'

			$scope.updateUsers()


		$scope.updateUsers = () ->
			offset = ($scope.pagination.pageNumber - 1) * $scope.pagination.perPage
			Admin.getUsers 'uid', offset, $scope.pagination.perPage
				.success (data) ->
					if data.count and data.users?
						$scope.users.splice 0
						data.users.forEach (u) -> $scope.users.push u
						# $scope.pagination.total = data.count
						$scope.pagination.total = 11

		$scope.isAdmin = () ->
			true # TODO: actually check whether current user is an admin

		$scope.showModal = (type, item) ->
			if item?
				$scope.modal.content.type = type
				$scope.modal.content.item = _.cloneDeep item
				$scope.modal.content.itemOriginal = item

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
			before = _.cloneDeep $scope.modal.content.itemOriginal
			after = _.cloneDeep $scope.modal.content.item

			after.usertype = parseInt after.usertype * 1

			diff = do ->
				result = {}
				Object.keys(before).forEach (k) ->
					if k is 'id' or before[k] isnt after[k]
						if k is 'reputation'
							result[k] = after[k] - before[k]
						else
							result[k] = after[k]
				result

			if Object.keys().length > 1
				# Make request
				console.log diff

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
			if newVal > 0 and newVal <= $scope.totalPages $scope.pagination.total
				$scope.pagination.pageNumber = newVal
				$scope.updateUsers()

		return
	]
