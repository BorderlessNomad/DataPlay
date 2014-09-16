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

    $scope.margin =
      top: 0
      right: 0
      bottom: 0
      left: 0

    $scope.mainsections = {
      d:
        title: 'Gov Departments/Boodies'
        colNameA: 'Entities'
        colNameB: 'Last 30 days'
        error: null
        graph: []
        items: []
      e:
        title: 'Political Events'
        colNameA: 'Event'
        colNameB: 'Last 30 days'
        error: null
        graph: []
        items: []
      r:
        title: 'Politically Aware/Active'
        colNameA: 'Location'
        colNameB: 'Last 30 days'
        error: null
        graph: []
        items: []
    }

    $scope.sidebarsections = []

    $scope.init = ->
      (['d', 'e', 'r']).forEach (i) ->
        OverviewScreen.get i
          .success (data) ->
            if data instanceof Array
              $scope.mainsections[i].items = data
              $scope.mainsections[i].items.forEach (item) ->
                total = 0
                for a in item.graph then total += a.y

                $scope.mainsections[i].graph.push
                  term: item.term
                  value: total

                item.id = "#{i.replace(/\W/g, '').toLowerCase()}-#{item.term.replace(/\W/g, '').toLowerCase()}"
                return
          .error $scope.handleError i

      OverviewScreen.get 'p'
        .success (data) ->
          if data instanceof Array
            $scope.sidebarsections = data.map (sect) ->
              sect.top5 = sect.top5.filter (item) ->
                return item.amount > 0
              return sect
        .error $scope.handleError 'p'

    $scope.renderLine = (details) ->
      return (chart) ->
        graph = details.graph

        entry = crossfilter graph
        dimension = entry.dimension (d) -> d.x
        group = dimension.group().reduceSum (d) -> d.y

        ordinals = []
        ordinals.push d.key for d in group.all() when d not in ordinals
        chart.colorAccessor (d, i) -> parseInt(d.y) % ordinals.length


        chart.dimension dimension
        chart.group group, "Test Title"

        xScale = d3.scale.linear()
          .domain d3.extent group.all(), (d) -> parseInt d.key
          .range [0, 60]
        chart.x xScale

        chart.keyAccessor (d) -> d.key
        chart.valueAccessor (d) -> d.value

        chart.xAxis().ticks 0
        chart.yAxis().ticks 0

        chart.xAxisLabel false, 0
        chart.yAxisLabel false, 0

        return

    $scope.renderPie = (details) ->
      return (chart) ->
        console.log "renderPie"
        console.log details
        graph = details.graph

        entry = crossfilter graph
        dimension = entry.dimension (d) -> d.term
        group = dimension.group().reduceSum (d) -> d.value

        chart.dimension dimension
        chart.group group

        chart.colorAccessor (d, i) -> i + 1
        chart.renderLabel false

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
