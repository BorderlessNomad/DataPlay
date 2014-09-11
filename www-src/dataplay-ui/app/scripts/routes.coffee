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
			.when '/search/:query',
				templateUrl: 'views/search.html'
				controller: 'SearchCtrl'
				login: true
			.when '/overview/:id',
				templateUrl: 'views/overview.html'
				controller: 'OverviewRelatedCtrl'
				login: true
			.when '/overview/:id/:offset/:count',
				templateUrl: 'views/overview.html'
				controller: 'OverviewRelatedCtrl'
				login: true
			.when '/charts/related/:id/:key/:type/:x/:y',
				templateUrl: 'views/charts.html'
				controller: 'ChartsRelatedCtrl'
				login: true
			.when '/charts/related/:id/:key/:type/:x/:y/:z',
				templateUrl: 'views/charts.html'
				controller: 'ChartsRelatedCtrl'
				login: true
			.when '/charts/correlated/:id/:correlationid/:type/:x/:y',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCorrelatedCtrl'
				login: true
			.when '/charts/correlated/:id/:correlationid/:type/:x/:y/:z',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCorrelatedCtrl'
				login: true
			.when '/map/:id',
				templateUrl: 'views/map.html'
				controller: 'MapCtrl'
				login: true
			.when '/user/login',
				templateUrl: 'views/user/login.html'
				login: false
			.when '/user/logout',
				templateUrl: 'views/user/login.html'
				controller: 'UserCtrl'
				login: true
			.when '/user/forgotpassword',
				templateUrl: 'views/user/forgot-password.html'
				controller: 'UserCtrl'
				login: false
			.when '/user/resetpassword/:token',
				templateUrl: 'views/user/reset-password.html'
				controller: 'UserCtrl'
				login: false
			.when '/user/register',
				templateUrl: 'views/user/register.html'
				login: false
			.otherwise
				redirectTo: '/'

		$locationProvider
			.html5Mode true
			.hashPrefix '!'

		return
	]

# Auth Handler
angular.module('dataplayApp')
	.run ['$rootScope', '$location', 'Auth', ($rootScope, $location, Auth) ->
		$rootScope.$on "$routeChangeStart", (event, nextRoute, currentRoute) ->
			if nextRoute? and nextRoute.login and not Auth.isAuthenticated()
				$location.path "/user/login"
				return

		return
	]

# Disable refresh on Route change when $location.path('/path', false)
angular.module('dataplayApp')
	.run ['$rootScope', '$location', '$route', ($rootScope, $location, $route) ->
		original = $location.path
		$location.path = (path, reload) ->
			if reload is false
				lastRoute = $route.current
				un = $rootScope.$on "$locationChangeSuccess", () ->
					$route.current = lastRoute
					un()

			original.apply $location, [path]

		return
	]
