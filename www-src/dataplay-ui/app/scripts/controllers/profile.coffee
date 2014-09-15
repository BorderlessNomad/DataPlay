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
    $scope.currentTab = if $scope.params?.tab?.length > 0 then $scope.params.tab else 'profile'
    $scope.error =
      message: null

    # Personal Details
    $scope.current =
      email: ''
      username: ''
    $scope.saved =
      email: ''
      username: ''
      success: false

    $scope.validDiscoveries = []
    $scope.discoveries = []
    $scope.observations = []

    $scope.inits =
      profile: ->
        Profile.getInfo()
          .success (res) ->
            $scope.current.email = res.data.email
            $scope.current.username = res.data.username

            $scope.saved.email = res.data.email
            $scope.saved.username = res.data.username
         .error (data, status) -> $scope.handleError data, status

      validdiscoveries: ->
        Profile.getValidDiscoveries()
          .success (data) ->
            if res instanceof Array
              $scope.validDiscoveries = res
            return
          .error (data, status) -> $scope.handleError data, status

      discoveries: ->
        Profile.getDiscoveries()
          .success (res) ->
            if res instanceof Array
              $scope.discoveries = res
            return
          .error (data, status) -> $scope.handleError data, status

      observations: ->
        Profile.getObservations()
          .success (res) ->
            if res instanceof Array
              $scope.observations = res
            return
          .error (data, status) -> $scope.handleError data, status

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

    $scope.hasError = () ->
      if $scope.error.message?.length then true else false

    $scope.handleError = (err, status) ->
      $scope.error.message = switch
        when err and err.message then err.message
        when err and err.data and err.data.message then err.data.message
        when err and err.data then err.data
        when err then err
        else ""

    return
  ]
