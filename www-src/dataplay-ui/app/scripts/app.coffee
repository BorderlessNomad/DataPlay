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
		'ipCookie'
		'ngResource'
		'ngRoute'
		'ngSanitize'
	])

angular.module('dataplayApp')
	.constant "config",
		cookieName: "DPSession"
		api:
			base_url: "http://localhost:3000/api"
