'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ProfileCtrl
 # @description
 # # ProfileCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'ProfileCtrl', ['$scope', '$location', '$routeParams', 'Profile', 'Auth', 'config', ($scope, $location, $routeParams, Profile, Auth, config) ->
    $scope.params = $routeParams
    $scope.currentTab = 'details'

    $scope.current =
      email: ''
      username: ''

    $scope.saved =
      email: ''
      username: ''

    $scope.inits =
      details: ->
        Profile.getInfo().then (res) ->
          $scope.current.email = res.data.email
          $scope.current.username = res.data.username
          $scope.saved.email = res.data.email
          $scope.saved.username = res.data.username

    $scope.inits[$scope.currentTab]?()

    $scope.changeTab = (tab) ->
      $scope.currentTab = tab
      $scope.inits[$scope.currentTab]?()

    $scope.submitDetails = ->
      Profile.setInfo $scope.current.email, $scope.current.username

    $scope.clearDetails = ->
      $scope.current.email = $scope.saved.email
      $scope.current.username = $scope.saved.username

    return
  ]
