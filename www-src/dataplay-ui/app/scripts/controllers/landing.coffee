'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:LandingCtrl
 # @description
 # # LandingCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'LandingCtrl', ['$scope', '$location', 'User', 'Overview', 'config', ($scope, $location, User, Overview, config) ->
    $scope.config = config

    $scope.stats =
      players: null
      discoveries: null
      datasets: null

    $scope.init = ->
      User.getStats()
        .success (data) ->
          if data instanceof Array
            data.forEach (d) ->
              $scope.stats[d.Label] = $scope.commarise d.Value

    $scope.commarise = (num) ->
      Number(num).toLocaleString()

    return
  ]
