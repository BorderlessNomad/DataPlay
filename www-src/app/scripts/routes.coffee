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
		$provide.decorator '$sniffer', ['$delegate', ($delegate) ->
			$delegate.history = false
			$delegate
		]

		$routeProvider
			.when '/home',
				templateUrl: 'views/home.html'
				controller: 'HomeCtrl'
				title: ['Home']
				login: true
			.when '/',
				templateUrl: 'views/landing.html'
				controller: 'LandingCtrl'
				login: false
			.when '/about',
				templateUrl: 'views/about.html'
				controller: 'LandingCtrl'
				title: ['About']
				login: false
			.when '/search',
				templateUrl: 'views/search.html'
				controller: 'SearchCtrl'
				title: ['Search']
				login: true
				preventReload: true
			.when '/search/:query',
				templateUrl: 'views/search.html'
				controller: 'SearchCtrl'
				title: ['Search']
				login: true
				preventReload: true
			.when '/overview',
				templateUrl: 'views/overviewscreen.html'
				controller: 'OverviewScreenCtrl'
				title: ['Overview']
				login: true
			.when '/overview/:id',
				templateUrl: 'views/overview.html'
				controller: 'OverviewCtrl'
				title: ['Overview']
				login: true
			.when '/overview/:id/:offset/:count',
				templateUrl: 'views/overview.html'
				controller: 'OverviewCtrl'
				title: ['Overview']
				login: true
			.when '/charts/related/:id/:key/:type/:x/:y',
				templateUrl: 'views/charts.html'
				controller: 'ChartsRelatedCtrl'
				title: ['Chart', 'Related']
				login: true
			.when '/charts/related/:id/:key/:type/:x/:y/:z',
				templateUrl: 'views/charts.html'
				controller: 'ChartsRelatedCtrl'
				title: ['Chart', 'Related']
				login: true
			.when '/charts/correlated/:id/:correlationid/:type/:x/:y',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCorrelatedCtrl'
				title: ['Chart', 'Correlated']
				login: true
			.when '/charts/correlated/:id/:correlationid/:type/:x/:y/:z',
				templateUrl: 'views/charts.html'
				controller: 'ChartsCorrelatedCtrl'
				title: ['Chart', 'Correlated']
				login: true
			.when '/map/:id',
				templateUrl: 'views/map.html'
				controller: 'MapCtrl'
				title: ['Map']
				login: true
			.when '/user/login',
				templateUrl: 'views/user/login.html'
				title: ['Login/Register']
				login: false
			.when '/user/logout',
				templateUrl: 'views/user/login.html'
				controller: 'UserCtrl'
				title: ['Logout']
				login: true
			.when '/user/forgotpassword',
				templateUrl: 'views/user/forgot-password.html'
				controller: 'UserCtrl'
				title: ['Forgot Password']
				login: false
			.when '/user/resetpassword/:token',
				templateUrl: 'views/user/reset-password.html'
				controller: 'UserCtrl'
				title: ['Reset Password']
				login: false
			.when '/user',
				templateUrl: 'views/user/profile.html'
				controller: 'ProfileCtrl'
				title: ['Profile']
				login: true
			.when '/user/:tab',# This must be kept after all /user/* calls to make sure that login and other known pages works.
				templateUrl: 'views/user/profile.html'
				controller: 'ProfileCtrl'
				title: ['Profile']
				login: true
			.when '/profile',
				redirectTo: '/user'
			.when '/profile/:user',
				templateUrl: 'views/user/profile.html'
				controller: 'ProfileCtrl'
				title: ['Profile']
				login: true
			.when '/profile/:user/:tab',
				templateUrl: 'views/user/profile.html'
				controller: 'ProfileCtrl'
				title: ['Profile']
				login: true
			.when '/admin',
				templateUrl: 'views/admin/dashboard.html'
				controller: 'AdminUsersCtrl'
				title: ['Admin']
				login: true
			.when '/admin/users',
				templateUrl: 'views/admin/users.html'
				controller: 'AdminUsersCtrl'
				title: ['Admin - Users']
				login: true
			.when '/admin/observations',
				templateUrl: 'views/admin/observations.html'
				controller: 'AdminObservationsCtrl'
				title: ['Admin - Observations']
				login: true
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

		$rootScope.$on "$routeChangeSuccess", (event, nextRoute, currentRoute) ->
			if nextRoute.hasOwnProperty("$$route") and nextRoute.$$route.title?.length > 0
				title = ""
				for val, key in nextRoute.$$route.title
					title += " - #{val}"
				$rootScope.title = title
				return

		return
	]

# Disable refresh on Route change when $location.path('/path', false)
# angular.module('dataplayApp')
# 	.run ['$rootScope', '$location', '$route', ($rootScope, $location, $route) ->
# 		original = $location.path
# 		$location.path = (path, reload) ->
# 			if reload is false
# 				lastRoute = $route.current
# 				un = $rootScope.$on "$locationChangeSuccess", () ->
# 					$route.current = lastRoute
# 					un()

# 			original.apply $location, [path]

# 		return
# 	]
