'use strict'

###*
 # @ngdoc directive
 # @name dataplayApp.directive:match
 # @description
 # # match
###
angular.module("dataplayApp")
	.directive "match", [() ->
		require: "ngModel"
		restrict: "A"
		scope:
			match: "="
		link: (scope, elem, attrs, ctrl) ->
			scope.$watch (->
				(ctrl.$pristine and angular.isUndefined(ctrl.$modelValue)) or scope.match is ctrl.$modelValue
			), (currentValue) ->
				ctrl.$setValidity "match", currentValue
				return

			return
	]
