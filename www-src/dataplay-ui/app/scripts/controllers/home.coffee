'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:HomeCtrl
 # @description
 # # HomeCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'HomeCtrl', ['$scope', '$location', 'User', 'Auth', 'Overview', 'config', ($scope, $location, User, Auth, Overview, config) ->
    $scope.config = config

    $scope.searchquery = ''

    $scope.validatePatterns = [
      {
        title: "A&E waiting times"
      }
      {
        title: "Crime Rate London"
      }
      {
        title: "GDP Prices"
      }
      {
        title: "Gold Prices"
      }
      {
        title: "NHS Spending"
      }
      {
        title: "Crime Rate London"
      }
    ]

    $scope.myActivity = [
      {
        text: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
      }
      {
        text: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
      }
    ]

    $scope.recentObservations = [
      {
        user:
          name: "Tom MySpace"
          avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        text: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
      }
      {
        user:
          name: "Tom MySpace"
          avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        text: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
      }
    ]

    $scope.dataExperts = [
      {
        rank: 1
        name: "Tom MySpace"
        avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        score: 205
        rankclass: 'gold'
      }
      {
        rank: 2
        name: "Tom MySpace"
        avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        score: 204
        rankclass: 'silver'
      }
      {
        rank: 3
        name: "Tom MySpace"
        avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        score: 203
        rankclass: 'bronze'
      }
      {
        rank: 4
        name: "Tom MySpace"
        avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        score: 202
      }
      {
        rank: 5
        name: "Tom MySpace"
        avatar: "https://pbs.twimg.com/profile_images/1237550450/mstom_400x400.jpg"
        score: 201
      }
    ]

    $scope.init = ->
      console.log "is home"

    $scope.search = ->
      $location.path "/search/#{$scope.searchquery}"

    return
  ]
