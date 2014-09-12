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

    # Personal Details
    $scope.current =
      email: ''
      username: ''
    $scope.saved =
      email: ''
      username: ''
      success: false

    # Valid Discoveries
    $scope.validDiscoveries = [
      {
        title: 'Gold Prices'
        link: 555
      }
      {
        title: 'A&E waiting times'
        title2: 'Crime Rate London'
        link: 556
        link2: 557
      }
      {
        title: 'GDP Prices'
        link: 558
      }
    ]

    # Discoveries
    $scope.discoveries = [
      {
        title: 'GDP Prices'
        link: 558
      }
    ]

    # Observations
    $scope.observations = [
      {
        title: 'Gold Prices'
        link: 555
        message: 'We should buy some gold!'
      }
      {
        title: 'A&E waiting times'
        title2: 'Crime Rate London'
        link: 556
        link2: 557
        message: 'I think this is really interesting!!'
      }
      {
        title: 'GDP Prices'
        link: 558
        message: 'Lets buy some EUR!'
      }
    ]

    $scope.inits =
      details: ->
        Profile.getInfo().then (res) ->
          $scope.current.email = res.data.email
          $scope.current.username = res.data.username

          $scope.saved.email = res.data.email
          $scope.saved.username = res.data.username

      validdiscoveries: ->
        Profile.getValidDiscoveries().then (res) ->
          console.log res

      discoveries: ->
        Profile.getDiscoveries().then (res) ->
          console.log res

      observations: ->
        Profile.getObservations().then (res) ->
          console.log res

    $scope.inits[$scope.currentTab]?()

    $scope.changeTab = (tab) ->
      $scope.currentTab = tab
      $scope.inits[$scope.currentTab]?()

    $scope.submitDetails = ->
      Profile.setInfo $scope.current.email, $scope.current.username
        .then (res) ->
          $scope.saved.email = $scope.current.email
          $scope.saved.username = $scope.current.username
          $scope.saved.success = true

    $scope.clearDetails = ->
      $scope.current.email = $scope.saved.email
      $scope.current.username = $scope.saved.username

    $scope.closeAlert = () ->
      $scope.saved.success = false

    return
  ]
