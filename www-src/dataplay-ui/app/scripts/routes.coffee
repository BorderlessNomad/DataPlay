'use strict'

###*
 # @ngdoc overview
 # @name dataplayApp
 # @description
 # # dataplayApp
 #
 # Routes for application.
###

angular.module('dataplayApp')
	.config ['$routeProvider', '$locationProvider', '$provide', ($routeProvider, $locationProvider, $provide) ->
		$provide.decorator '$sniffer', ($delegate) ->
			$delegate.history = false
			$delegate

		$routeProvider
			.when '/',
				templateUrl: 'views/home.html'
				controller: 'HomeCtrl'
				login: true
			.when '/search',
				templateUrl: 'views/search.html'
				controller: 'SearchCtrl'
				login: true
			.when '/overview',
				templateUrl: 'views/overview.html'
				controller: 'OverviewCtrl'
				login: true
			.when '/charts',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCtrl'
				login: true
			.when '/login',
				templateUrl: 'views/login.html'
				login: false
			.when '/logout',
				templateUrl: 'views/login.html'
				controller: 'UserCtrl'
				login: true
			.when '/register',
				templateUrl: 'views/register.html'
				login: false
			.otherwise
				redirectTo: '/'

		$locationProvider
			.html5Mode true
			.hashPrefix '!'

		return
	]

angular.module('dataplayApp')
	.run ['$rootScope', '$location', 'Auth', ($rootScope, $location, Auth) ->
		$rootScope.$on "$routeChangeStart", (event, nextRoute, currentRoute) ->
			if nextRoute? and nextRoute.login and Auth.isAuthenticated() is false
				$location.path "/login"
				return

		return
	]
