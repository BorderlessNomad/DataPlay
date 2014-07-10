'use strict'

###*
 # @ngdoc overview
 # @name dataplayApp
 # @description
 # # dataplayApp
 #
 # Main module of the application.
###
angular
	.module('dataplayApp', [
		'ngAnimate'
		'ngCookies'
		'ngResource'
		'ngRoute'
		'ngSanitize'
		'UserApp'
	])

angular.module('dataplayApp')
	.constant "config",
		api:
			base_url: "http://localhost:3000/api"
