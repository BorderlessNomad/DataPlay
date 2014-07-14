'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCtrl
 # @description
 # # ChartsCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCtrl', ['$scope', ($scope) ->
		$scope.awesomeThings = [
			'HTML5 Boilerplate'
			'AngularJS'
			'Karma'
		]
	]
