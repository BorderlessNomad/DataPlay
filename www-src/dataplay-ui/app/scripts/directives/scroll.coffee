'use strict'

###*
 # @ngdoc directive
 # @name dataplayApp.directive:scroll
 # @description
 # # scroll
###
fakeNgModel = (initValue) ->
	$setViewValue: (value) ->
		@$viewValue = value
		return

angular.module('dataplayApp')
	.directive "scroll-glue", [() ->
		priority: 1
		require: ['?ngModel']
		restrict: "A"
		link: (scope, elem, attrs, ctrl) ->
			el = elem[0]
			ngModel = ctrl[0] || fakeNgModel true

			scrollToBottom = ->
				el.scrollTop = el.scrollHeight
				return

			shouldActivateAutoScroll = ->
				el.scrollTop + el.clientHeight + 1 >= el.scrollHeight

			scope.$watch ->
				scrollToBottom() if ngModel.$viewValue
				return

			elem.bind "scroll", ->
				activate = shouldActivateAutoScroll
				$scope.$apply ngModel.$setViewValue.bind(ngModel, activate) if activate isnt ngModel.$viewValue
				return

			return
	]
