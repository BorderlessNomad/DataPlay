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

		$scope.pagination =
			orderby: 'uid'
			perPage: 10
			pageNumber: 1
			total: 0

		$scope.users = []


		# General controls
		$scope.init = () ->
			if not $scope.isAdmin()
				$location.path '/home'

			$scope.updateUsers()


		$scope.updateUsers = (cb = (() ->)) ->
			offset = ($scope.pagination.pageNumber - 1) * $scope.pagination.perPage
			Admin.getUsers $scope.pagination.orderby, offset, $scope.pagination.perPage
				.success (data) ->
					if data.count and data.users?
						$scope.users.splice 0
						data.users.forEach (u) ->
							$scope.users.push
								uid: u.uid || 0
								avatar: u.avatar || ''
								username: u.username || ''
								email: u.email || ''
								md5email: u.md5email || ''
								reputation: u.reputation || 0
								usertype: u.usertype || 0
								enabled: if not u.enabled? then true else u.enabled
								password: ''
						$scope.pagination.total = data.count
					cb()

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

		$scope.generateRandom = () ->
			$scope.modal.content.item.password = do ->
				chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345679-_.'
				result = ''
				for i in [0..(12 + Math.floor(Math.random() * 8))]
					result += chars.charAt Math.floor(Math.random() * chars.length)
				result

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
					if k is 'uid' or before[k] isnt after[k]
						if k is 'reputation'
							result[k] = after[k] - before[k]
						else
							result[k] = after[k]
				result

			if Object.keys(diff).length > 1
				Admin.editUser diff
					.success (data) ->
						$scope.updateUsers () ->
							$scope.closeModal()
			return

		$scope.disable = () ->
			params =
				uid: $scope.modal.content.item.uid
				enabled: not $scope.modal.content.item.enabled
			Admin.editUser params
				.success (data) ->
					$scope.updateUsers () ->
						$scope.closeModal()
			return



		# Pagination
		$scope.orderby = (col) ->
			if $scope.pagination.orderby isnt col
				$scope.pagination.orderby = col
				$scope.pagination.pageNumber = 1
				$scope.updateUsers()

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
