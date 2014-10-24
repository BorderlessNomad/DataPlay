'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ProfileCtrl
 # @description
 # # ProfileCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ProfileCtrl', ['$scope', '$location', '$routeParams', 'Profile', 'Auth', 'config', ($scope, $location, $routeParams, Profile, Auth, config) ->
		$scope.params = $routeParams
		$scope.currentTab = if $scope.params?.tab?.length > 0 then $scope.params.tab else 'profile'
		$scope.error =
			message: null

		$scope.loading =
			profile: false
			creditdiscoveries: false
			discoveries: false
			observations: false

		# Personal Details
		$scope.current =
			email: ''
			username: ''
		$scope.saved =
			email: ''
			username: ''
			success: false

		$scope.creditDiscoveries = []
		$scope.discoveries = []
		$scope.observations = []

		$scope.inits =
			profile: ->
				$scope.loading.profile = true
				Profile.getInfo()
					.success (data) ->
						$scope.loading.profile = false
						$scope.current.email = data.email
						$scope.current.username = data.username

						$scope.saved.email = data.email
						$scope.saved.username = data.username
					.error (data, status) ->
						$scope.loading.profile = false
						$scope.handleError data, status

			creditdiscoveries: ->
				$scope.loading.creditdiscoveries = true
				Profile.getCreditDiscoveries()
					.success (data) ->
						$scope.loading.creditdiscoveries = false
						if data instanceof Array
							$scope.creditDiscoveries = data
						return
					.error (data, status) ->
						$scope.loading.creditdiscoveries = false
						$scope.handleError data, status

			discoveries: ->
				$scope.loading.discoveries = true
				Profile.getDiscoveries()
					.success (data) ->
						$scope.loading.discoveries = false
						if data instanceof Array
							$scope.discoveries = data
						return
					.error (data, status) ->
						$scope.loading.discoveries = false
						$scope.handleError data, status

			observations: ->
				$scope.loading.observations = true
				Profile.getObservations()
					.success (data) ->
						$scope.loading.observations = false
						if data instanceof Array
							$scope.observations = data
						return
					.error (data, status) ->
						$scope.loading.observations = false
						$scope.handleError data, status

		$scope.inits[$scope.currentTab]?()

		$scope.changeTab = (tab) ->
			$location.path "/user/#{tab}"

		$scope.submitDetails = ->
			Profile.setInfo $scope.current.email, $scope.current.username
				.then (res) ->
					$scope.saved.email = $scope.current.email
					$scope.saved.username = $scope.current.username
					$scope.saved.success = true

		$scope.clearDetails = ->
			$scope.current.email = $scope.saved.email
			$scope.current.username = $scope.saved.username

		$scope.closeAlert = () ->
			$scope.saved.success = false

		$scope.hasError = () ->
			if $scope.error.message?.length then true else false

		$scope.handleError = (err, status) ->
			$scope.error.message = switch
				when err and err.message then err.message
				when err and err.data and err.data.message then err.data.message
				when err and err.data then err.data
				when err then err
				else ""

			if $scope.error.message.substring(0, 6) is '<html>'
				$scope.error.message = do ->
					curr = $scope.error.message
					curr = curr.replace(/(\r\n|\n|\r)/gm, '')
					curr = curr.replace(/.{0,}(\<title\>)/, '')
					curr = curr.replace(/(\<\/title\>).{0,}/, '')
					curr

		return
	]
