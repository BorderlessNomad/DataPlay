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
		$scope.currentTab = if $scope.currentTab is 'profile' and $scope.params?.user? then 'approveddiscoveries' else $scope.currentTab
		$scope.error =
			message: null

		$scope.loading =
			profile: false
			approveddiscoveries: false
			discoveries: false
			observations: false

		# Personal Details
		$scope.current =
			email: ''
			username: ''
			avatar: ''
		$scope.saved =
			email: ''
			username: ''
			success: false

		$scope.isLoggedInUser = /^(\/user)/.test $location.url()

		if not $scope.isLoggedInUser
			$scope.current.username = $routeParams.user
			loggedInUsername = Auth.get config.userName
			if $scope.current.username is loggedInUsername
				$location.path('/user/profile').replace()

		$scope.approvedDiscoveries = []
		$scope.discoveries = []
		$scope.observations = []

		$scope.profile = (user = '') ->
			$scope.loading.profile = true

			Profile.getInfo(user)
				.success (data) ->
					$scope.loading.profile = false
					$scope.current.email = data.email
					$scope.current.username = data.username
					$scope.current.avatar = data.avatar or "http://www.gravatar.com/avatar/#{data.email_hash}?d=identicon"

					if not user? or not user.length
						$scope.saved.email = data.email
						$scope.saved.username = data.username

					return
				.error (data, status) ->
					$scope.loading.profile = false
					$scope.handleError data, status

		$scope.approveddiscoveries = (user = '') ->
			$scope.loading.approveddiscoveries = true

			$scope.profile user

			Profile.getCreditDiscoveries(user)
				.success (data) ->
					$scope.loading.approveddiscoveries = false
					if data instanceof Array
						$scope.approvedDiscoveries = data

					return
				.error (data, status) ->
					$scope.loading.approveddiscoveries = false
					$scope.handleError data, status

		$scope.discoveries = (user = '') ->
			$scope.loading.discoveries = true

			$scope.profile user

			Profile.getDiscoveries(user)
				.success (data) ->
					$scope.loading.discoveries = false
					if data instanceof Array
						$scope.discoveries = data

					return
				.error (data, status) ->
					$scope.loading.discoveries = false
					$scope.handleError data, status

		$scope.observations = (user = '') ->
			$scope.loading.observations = true

			$scope.profile user

			Profile.getObservations(user)
				.success (data) ->
					$scope.loading.observations = false
					if data instanceof Array
						$scope.observations = data

					return
				.error (data, status) ->
					$scope.loading.observations = false
					$scope.handleError data, status

		$scope[$scope.currentTab]?($scope.params.user ? null)

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
