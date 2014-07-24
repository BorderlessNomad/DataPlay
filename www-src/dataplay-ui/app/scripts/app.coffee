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
		'ipCookie'
		'ui.bootstrap'
	])

angular.module('dataplayApp')
	.constant "config",
		sessionHeader: "X-API-SESSION"
		sessionName: "DPSession"
		userName: "DPUser"
		api:
			base_url: "http://localhost:3000/api"
