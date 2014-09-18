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
          console.log data
          # $scope.stats.players = "200,000"
          # $scope.stats.discoveries = "200,000"
          # $scope.stats.datasets = "200,000"

    $scope.commarise = (num) ->
      Number(num).toLocaleString()

    return
  ]
