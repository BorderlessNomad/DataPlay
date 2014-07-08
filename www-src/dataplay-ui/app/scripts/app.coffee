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
	.config ($routeProvider) ->
		$routeProvider
			.when '/',
				templateUrl: 'views/home.html'
				controller: 'HomeCtrl'
			.when '/search',
				templateUrl: 'views/search.html'
				controller: 'SearchCtrl'
			.when '/overview',
				templateUrl: 'views/overview.html'
				controller: 'OverviewCtrl'
			.when '/charts',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCtrl'
			.when '/login',
				templateUrl: 'views/login.html'
				login: true
			.when '/logout',
				templateUrl: 'views/login.html'
			.when '/register',
				templateUrl: 'views/register.html'
				public: true
			.otherwise
				redirectTo: '/'

angular.module('dataplayApp')
	.run (user) ->
		user.init
			appId: "dataplayApp"
