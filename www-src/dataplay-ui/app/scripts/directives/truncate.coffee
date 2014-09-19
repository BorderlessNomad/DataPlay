'use strict'

###*
 # @ngdoc directive
 # @name dataplayApp.directive:truncate
 # @description
 # # truncate
 #
 # Usage:
 #	{{some_text | truncate:25:true:' ..'}}
 # Options:
 #	max (integer, default: 10) - max length of the text, cut to this number of chars.
 #	wordwise (boolean, default: false) - if true, cut only by words bounds.
 #	tail (string, default: " ...") - add this string to the input string if the string was cut.
###
angular.module("dataplayApp")
	.filter "truncate", [() ->
		(value, max = 10, wordwise = false, tail = " ...") ->
			return "" if not value?

			return value if value.length <= max

			value = value.substr 0, max

			if wordwise
				boundary = value.lastIndexOf " "
				value = value.substr(0, boundary) unless boundary is -1

			"#{value}#{tail}"
	]
