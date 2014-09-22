'use strict'

###*
 # @ngdoc service
 # @name dataplayApp.OverviewScreen
 # @description
 # # OverviewScreen
 # Factory in the dataplayApp.
###
angular.module('dataplayApp')
  .factory 'OverviewScreen', ['$http', 'Auth', 'config', ($http, Auth, config) ->
    get: (type = 'd') ->
      $http.get config.api.base_url + "/political/#{type}"
  ]
