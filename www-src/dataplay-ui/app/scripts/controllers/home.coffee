'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:HomeCtrl
 # @description
 # # HomeCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'HomeCtrl', ['$scope', '$location', 'Home', 'Auth', 'Overview', 'config', ($scope, $location, Home, Auth, Overview, config) ->
    $scope.config = config

    $scope.searchquery = ''

    $scope.validatePatterns = null

    $scope.myActivity = null
    $scope.recentObservations = null
    $scope.dataExperts = null

    $scope.init = ->
      Home.getAwaitingValidation()
        .success (data) ->
          $scope.validatePatterns = []
          console.log data
        .error ->
          $scope.validatePatterns = []


      Home.getActivityStream()
        .success (data) ->
          if data instanceof Array
            $scope.myActivity = data.map (d) ->
              date: Overview.humanDate new Date d.time
              pretext: d.activitystring
              linktext: d.patternid
              url: d.linkstring
          else
            $scope.myActivity = []
        .error ->
          $scope.myActivity = []

      Home.getRecentObservations()
        .success (data) ->
          if data instanceof Array
            $scope.recentObservations = data.map (d) ->
              user:
                name: d.username
                avatar: "http://www.gravatar.com/avatar/#{d.MD5email}?d=identicon"
              text: d.comment
              url: d.linkstring
          else
            $scope.recentObservations = []
        .error ->
          $scope.recentObservations = []

      Home.getDataExperts()
        .success (data) ->
          if data instanceof Array

            medals = ['gold', 'silver', 'bronze']

            $scope.dataExperts = data.map (d, key) ->
              obj =
                rank: key + 1
                name: d.username
                avatar: "http://www.gravatar.com/avatar/#{d.MD5email}?d=identicon"
                score: d.reputation

              if obj.rank <= 3 then obj.rankclass = medals[obj.rank - 1]

              obj
          else
            $scope.dataExperts = []
        .error ->
          $scope.dataExperts = []

    $scope.search = ->
      $location.path "/search/#{$scope.searchquery}"

    return
  ]
