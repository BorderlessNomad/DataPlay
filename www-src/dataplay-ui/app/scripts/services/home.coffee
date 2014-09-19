'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.Home
 # @description
 # # Home
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
  .factory 'Home', ['$http', 'Auth', 'config', ($http, Auth, config) ->

    getStats: () ->
      $http.get config.api.base_url + "/home/data"

    getActivityStream: () ->
      $http.get config.api.base_url + "/user/activitystream"

    getRecentObservations: () ->
      $http.get config.api.base_url + "/recentobservations"

    getDataExperts: () ->
      $http.get config.api.base_url + "/user/experts"
  ]
