'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewCtrl
 # @description
 # # OverviewCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'OverviewCtrl', ['$scope', '$routeParams', '$location', 'Overview', ($scope, $routeParams, $location, Overview) ->
		$scope.params = $routeParams
		$scope.info =
			name: $scope.params.id
			title: $scope.params.id

		$scope.error = null

		$scope.info = () ->
			Overview.info $scope.params.id
				.success (data) ->
					$scope.info = data

					if $scope.info.name isnt $scope.params.id
						$location.path("/overview/#{$scope.info.name}").replace()

					console.log $scope.info
					return
				.error $scope.handleError

		$scope.handleError = (err, status) ->
			console.log "Overview::info::Error:", status

			$scope.error = switch
				when err and err.message then err.message
				when err and err.data and err.data.message then err.data.message
				when err and err.data then err.data
				when err then err
				else ''

			if $scope.error.substring(0, 6) is '<html>'
				$scope.error = do ->
					curr = $scope.error
					curr = curr.replace(/(\r\n|\n|\r)/gm, '')
					curr = curr.replace(/.{0,}(\<title\>)/, '')
					curr = curr.replace(/(\<\/title\>).{0,}/, '')
					curr
	]
