'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:LandingCtrl
 # @description
 # # LandingCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
  .controller 'LandingCtrl', ['$scope', '$location', 'Home', 'Auth', 'Overview', 'PatternMatcher', 'config', ($scope, $location, Home, Auth, Overview, PatternMatcher, config) ->
    $scope.config = config
    $scope.Auth = Auth
    $scope.username = Auth.get config.userName

    $scope.chartsRelated = []

    $scope.relatedChart = new RelatedCharts $scope.chartsRelated

    # $scope.relatedChart.width = 200
    # $scope.relatedChart.height = 150

    if Auth.isAuthenticated()
      $location.path '/home'

    $scope.stats =
      players: null
      discoveries: null
      datasets: null

    $scope.init = ->
      Home.getStats()
        .success (data) ->
          if data instanceof Array
            data.forEach (d) ->
              $scope.stats[d.Label] = $scope.commarise d.Value

      Home.getTopRated()
        .success (data) ->
          if data? and data.length > 0
            for key, chart of data
              continue unless $scope.relatedChart.isPlotAllowed chart.type
              continue unless chart.relationid?

              guid = chart.relationid.split("/")[0]

              key = parseInt(key)
              chart.key = key
              chart.id = "related-#{guid}-#{chart.key + $scope.relatedChart.offset.related}-#{chart.type}"
              chart.url = "charts/related/#{guid}/#{chart.key}/#{chart.type}/#{chart.xLabel}/#{chart.yLabel}"
              chart.url += "/#{chart.zLabel}" if chart.type is 'bubble'

              chart.patterns = {}
              chart.patterns[chart.xLabel] =
                valuePattern: PatternMatcher.getPattern chart.values[0]['x']
                keyPattern: PatternMatcher.getKeyPattern chart.values[0]['x']

              if chart.patterns[chart.xLabel].valuePattern is 'date'
                for value, key in chart.values
                  chart.values[key].x = new Date(value.x)

              if chart.yLabel?
                chart.patterns[chart.yLabel] =
                  valuePattern: PatternMatcher.getPattern chart.values[0]['y']
                  keyPattern: PatternMatcher.getKeyPattern chart.values[0]['y']

              $scope.chartsRelated.push chart if PatternMatcher.includePattern(
                chart.patterns[chart.xLabel].valuePattern,
                chart.patterns[chart.xLabel].keyPattern
              )

    $scope.commarise = (num) ->
      Number(num).toLocaleString()

    $scope.width = $scope.relatedChart.width
    $scope.height = $scope.relatedChart.height
    $scope.margin = $scope.relatedChart.margin

    $scope.hasRelatedCharts = $scope.relatedChart.hasRelatedCharts
    $scope.lineChartPostSetup = $scope.relatedChart.lineChartPostSetup
    $scope.rowChartPostSetup = $scope.relatedChart.rowChartPostSetup
    $scope.columnChartPostSetup = $scope.relatedChart.columnChartPostSetup
    $scope.pieChartPostSetup = $scope.relatedChart.pieChartPostSetup
    $scope.bubbleChartPostSetup = $scope.relatedChart.bubbleChartPostSetup

    return
  ]
