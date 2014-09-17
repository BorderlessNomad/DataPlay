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

		$scope.mainSections = {
			d:
				title: 'Gov Departments/Bodies'
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

		$scope.sidebarSections = []

		$scope.init = ->
			(['d', 'e', 'r']).forEach (i) ->
				OverviewScreen.get i
					.success (data) ->
						if data instanceof Array
							$scope.mainSections[i].items = data
							$scope.mainSections[i].items.forEach (item) ->
								total = 0
								for a in item.graph then total += a.y

								$scope.mainSections[i].graph.push
									term: item.term
									value: total

								item.id = "#{i.replace(/\W/g, '').toLowerCase()}-#{item.term.replace(/\W/g, '').toLowerCase()}"

								return
					.error $scope.handleError i

			OverviewScreen.get 'p'
				.success (data) ->
					if data instanceof Array
						$scope.sidebarSections = data.map (sect) ->
							sect.url = if sect.id is "top_discoverers" then "user/profile" else "search"
							sect.top5 = sect.top5.filter (item) ->
								item.amount > 0

							sect
				.error $scope.handleError 'p'

		$scope.renderLine = (details) ->
			(chart) ->
				graph = details.graph

				entry = crossfilter graph
				dimension = entry.dimension (d) -> d.x
				group = dimension.group().reduceSum (d) -> d.y

				ordinals = []
				ordinals.push d.key for d in group.all() when d not in ordinals
				chart.colorAccessor (d, i) -> parseInt(d.y) % ordinals.length

				chart.dimension dimension
				chart.group group

				xScale = d3.scale.linear()
					.domain d3.extent group.all(), (d) -> parseInt d.key
					.range [0, 60]
				chart.x xScale

				chart.keyAccessor (d) -> d.key
				chart.valueAccessor (d) -> d.value

				chart.xAxis().ticks(0).tickFormat (v) -> ""
				chart.yAxis().ticks(0).tickFormat (v) -> ""

				chart.xAxisLabel false, 0
				chart.yAxisLabel false, 0

				return

		$scope.renderPie = (details) ->
			(chart) ->
				graph = details.graph

				entry = crossfilter graph
				dimension = entry.dimension (d) -> d.term

				groupSum = 0
				group = dimension.group().reduceSum (d) ->
					groupSum += d.value
					d.value

				chart.dimension dimension
				chart.group group

				chart.colorAccessor (d, i) -> i + 1

				chart.renderlet (c) ->
					svg = d3.select 'svg'
					slices = c.svg().selectAll "g.pie-slice"
					slices.on "mouseover", (d) ->
						percent = (d.data.value / groupSum * 100).toFixed(2)
						d3.select "#pie-tooltip"
							.style "left", "#{d3.event.pageX - 400}px"
							.style "top", "#{d3.event.pageY - 55}px"
							.attr "class", "tooltip fade top in"
							.select ".tooltip-inner"
								.text "#{d.data.key}: #{d.data.value} (#{percent}%)"
					slices.on "mouseout", (d) ->
						d3.select "#pie-tooltip"
							.attr "class", "tooltip top hidden"

		$scope.handleError = (type) ->
			(err, status) ->
				$scope.mainSections[type].error = switch
					when err and err.message then err.message
					when err and err.data and err.data.message then err.data.message
					when err and err.data then err.data
					when err then err
					else ''

		return
	]
