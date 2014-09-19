'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:LandingCtrl
 # @description
 # # LandingCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'LandingCtrl', ['$scope', '$location', 'Home', 'Auth', 'Overview', 'config', ($scope, $location, Home, Auth, Overview, config) ->
    $scope.config = config
    $scope.Auth = Auth
    $scope.username = Auth.get config.userName

    if Auth.isAuthenticated()
      $location.path '/home'

    $scope.stats =
      players: null
      discoveries: null
      datasets: null

    $scope.init = ->
      Home.getStats()
        .success (data) ->
          if data instanceof Array
            data.forEach (d) ->
              $scope.stats[d.Label] = $scope.commarise d.Value

    $scope.commarise = (num) ->
      Number(num).toLocaleString()

    return
  ]
