'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:OverviewScreenCtrl
 # @description
 # # OverviewScreenCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'OverviewScreenCtrl', ['$scope', '$location', '$routeParams', 'OverviewScreen', 'Auth', 'config', ($scope, $location, $routeParams, OverviewScreen, Auth, config) ->
    $scope.params = $routeParams

    $scope.mainsections = {
      d:
        title: 'Gov Departments/Boodies'
        colNameA: 'Entities'
        colNameB: 'Last 30 days'
        error: null
        items: [
          {
            title: 'Health'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
          {
            title: 'Cabinet Office'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
        ]
      e:
        title: 'Political Events'
        colNameA: 'Event'
        colNameB: 'Last 30 days'
        error: null
        items: [
          {
            title: 'Announcement'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
          {
            title: 'Protest'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
        ]
      r:
        title: 'Politically Aware/Active'
        colNameA: 'Location'
        colNameB: 'Last 30 days'
        error: null
        items: [
          {
            title: 'Westminster'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
          {
            title: 'Westminster'
            data: [0,1,2,3,4,5,4,3,2,1]
          }
        ]
    }

    $scope.sidebarsections = [
      {
        title: 'Most Popular keywords'
        items: [
          {
            title: 'A&E'
            count: 22241
          }
          {
            title: 'London'
            count: 22240
          }
          {
            title: 'Alcohol'
            count: 22239
          }
        ]
      }
      {
        title: 'Top Correlated keywords'
        items: [
          {
            title: 'A&E'
            count: 23241
          }
          {
            title: 'London'
            count: 23240
          }
          {
            title: 'Alcohol'
            count: 23239
          }
        ]
      }
      {
        title: 'Top Discoverers'
        items: []
      }
    ]

    $scope.init = ->
      (['d', 'e', 'r']).forEach (i) ->
        OverviewScreen.get i
          .success (data) ->
            console.log data
          .error $scope.handleError i

    $scope.handleError = (type) ->
      return (err, status) ->
        $scope.mainsections[type].error = switch
          when err and err.message then err.message
          when err and err.data and err.data.message then err.data.message
          when err and err.data then err.data
          when err then err
          else ''

    $scope.init()
    return
  ]
