'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:AdminObservationsCtrl
 # @description
 # # AdminObservationsCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'AdminObservationsCtrl', ['$scope', '$location', 'Admin', 'Auth', 'Overview', 'config', ($scope, $location, Admin, Auth, Overview, config) ->

		$scope.modal =
			shown: false
			content:
				title: null
				type: null
				item: null

		$scope.pagination =
			flaggedonly: false
			orderby: 'observation_id'
			perPage: 10
			pageNumber: 1
			total: 0

		$scope.observations = []
		$scope.loading = false


		# General controls
		$scope.init = () ->
			$scope.bootNonAdmins()
			$scope.updateObservations()

		$scope.bootNonAdmins = () ->
			if not Auth.isAdmin()
				$location.path '/home'

		$scope.updateObservations = (cb = (() ->)) ->
			offset = ($scope.pagination.pageNumber - 1) * $scope.pagination.perPage
			$scope.observations.splice 0
			$scope.loading = true
			Admin.getObservations $scope.pagination.orderby, offset, $scope.pagination.perPage, $scope.pagination.flaggedonly
				.success (data) ->
					$scope.loading = false
					if data.count and data.comments?
						data.comments.forEach (u) ->
							$scope.observations.push
								observationid: u.observationid || 0
								created: u.created
								comment: u.comment || 0
								uid: u.uid || 0
								username: u.username || ''
								rating: (u.credited || 0) - (u.discredited)
								flagged: if not u.flagged? then false else u.flagged
						$scope.pagination.total = data.count

						# if there's no items on this page, but there are items overall, automatically revert to page 1
						if $scope.pagination.pageNumber isnt 1 &&
								$scope.pagination.total isnt 0 &&
								$scope.observations.length is 0
							$scope.pagination.pageNumber = 1
							$scope.updateObservations cb
							return
					cb()
				.error ->
					$scope.loading = false

		$scope.humanDate = (str) ->
			days = ['Sun', 'Mon', 'Tues', 'Wed', 'Thu', 'Fri', 'Sat']
			monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
			dat = new Date str
			if dat.toString() is 'Invalid Date'
				dat = new Date 0
			"#{days[dat.getDay()]} #{dat.getDate()} #{monthNames[dat.getMonth()]}, #{dat.getFullYear()}"


		$scope.showModal = (type, item) ->
			if item?
				$scope.modal.content.type = type
				$scope.modal.content.item = _.cloneDeep item
				$scope.modal.content.itemOriginal = item

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

		$scope.delete = () ->
			Admin.deleteObservation $scope.modal.content.item.observationid
				.success (data) ->
					$scope.updateObservations () ->
						$scope.closeModal()
			return



		# Pagination
		$scope.orderby = (col) ->
			if $scope.pagination.orderby isnt col
				$scope.pagination.orderby = col
				$scope.pagination.pageNumber = 1
				$scope.updateObservations()

		$scope.totalPages = (total) ->
			Math.ceil total / $scope.pagination.perPage

		$scope.range = (usrsLen) ->
			end = $scope.totalPages usrsLen
			if end is 0 then end = 1
			[1..end]

		$scope.changePage = (page, add) ->
			newVal = if add? then $scope.pagination.pageNumber + page else page
			if newVal > 0 and newVal <= $scope.totalPages $scope.pagination.total
				$scope.pagination.pageNumber = newVal
				$scope.updateObservations()

		$scope.setFilter = (flaggedonly) ->
			$scope.pagination.flaggedonly = flaggedonly
			$scope.updateObservations()

		return
	]
