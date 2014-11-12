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
		$scope.mapGen = new MapGenerator '#regionMap'

		$scope.margin =
			top: 0
			right: 0
			bottom: 0
			left: 0

		$scope.mainSections =
			d:
				title: 'Gov Departments/Bodies'
				colNameA: 'Entities'
				colNameB: 'Last 30 days'
				error: null
				type: 'pie'
				graph: []
				items: []
			e:
				title: 'Political Events'
				colNameA: 'Event'
				colNameB: 'Last 30 days'
				error: null
				type: 'pie'
				graph: []
				items: []
			r:
				title: 'Politically Aware/Active'
				colNameA: 'Location'
				colNameB: 'Last 30 days'
				error: null
				type: 'map'
				graph: []
				items: []

		$scope.colourMap = {}

		$scope.sidebarSections = []

		$scope.init = ->
			(['d', 'e', 'r']).forEach (i) ->
				OverviewScreen.get i
					.success (data) ->
						if data instanceof Array
							$scope.mainSections[i].items = data.filter (i) ->
								!! i.term
							maxTotal = 0

							$scope.mainSections[i].items.forEach (item) ->
								total = 0
								for a in item.graph then total += a.y

								if total > maxTotal then maxTotal = total

								item.slug = item.term.toLowerCase().replace(/_|\-|\'|\s/g, '')
								item.id = "#{i.replace(/\W/g, '').toLowerCase()}-#{item.slug}"

								$scope.mainSections[i].graph.push
									id: item.id
									slug: item.slug
									term: item.term
									value: total

								return

							if $scope.mainSections[i].type is 'map'
								dictionary = {}
								$scope.mainSections[i].graph.forEach (item) ->
									dictionary[item.term] = item

								$scope.mapGen.maxvalue = maxTotal

								regData = Object.keys($scope.mapGen.boundaryPaths).map (c) ->
									name: c
									value: dictionary[c]?.value || 0

								$scope.mapGen.generate regData

								$scope.mainSections[i].items.forEach (item) ->
									item.corresponds = "region-#{item.slug}"
									item.color = $scope.mapGen.getColor dictionary[item.term]?.value || 0

								regions = d3.selectAll '.region'

								regions.on "mouseover", (d) ->
									x = d3.event.pageX - $(window.document).scrollLeft()
									y = d3.event.pageY - $(window.document).scrollTop()

									el = d3.select @

									region = _.find $scope.mainSections[i].graph, (it) ->
										return it.id.replace('r-', '') is el.attr('id').replace('region-', '')
									if not region?
										region =
											id: "r-#{el.attr('id').replace('region-', '')}"
											slug: el.attr('id').replace('region-', '')
											term: el.attr 'data-display'
											value: 0

									d3.select "#pie-tooltip"
										.style "left", "#{x}px"
										.style "top", "#{y}px"
										.attr "class", "tooltip fade top in"
										.select ".tooltip-inner"
											.text "#{region.term}: #{region.value}"

								regions.on "mouseout", (d) ->
									d3.select "#pie-tooltip"
										.attr "class", "tooltip top hidden"

					.error $scope.handleError i

			OverviewScreen.get 'p'
				.success (data) ->
					if data instanceof Array
						$scope.sidebarSections = data.map (sect) ->
							sect.url = if sect.id is "top_discoverers" then "user/profile" else "search"
							sect.top = sect.top.filter (item) ->
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

				dimensionMap = {}
				dimension = entry.dimension (d) ->
					dimensionMap["#{d.term}-#{d.value}"] = d.id
					d.term

				groupSum = 0
				group = dimension.group().reduceSum (d) ->
					groupSum += d.value
					d.value

				chart.dimension dimension
				chart.group group

				dataRange = do ->
					curr =
						min: 100
						max: 0
					details.graph.forEach (g) ->
						if g.value > curr.max then curr.max = g.value
						if g.value < curr.min then curr.min = g.value
					curr

				chart.colors (d, i, j) ->
					curr = 0
					details.graph.forEach (g) ->
						if g.term is d then curr = g.value
					$scope.mapGen.getColor curr - dataRange.min, dataRange.max - dataRange.min

				chart.renderlet (c) ->
					svg = d3.select 'svg'
					slices = c.svg().selectAll "g.pie-slice"

					slices.each (d) ->
						id = dimensionMap["#{d.data.key}-#{d.data.value}"]
						slice = d3.select @
						slice.attr "id", "slice-#{id}"

						color = slice.select("path").attr("fill")

						$scope.colourMap[id] = color
						# legend = angular.element document.querySelector "legend-#{id}"

					slices.on "mouseover", (d) ->
						slice = d3.select @

						x = d3.event.pageX - $(window.document).scrollLeft()
						y = d3.event.pageY - $(window.document).scrollTop()

						percent = (d.data.value / groupSum * 100).toFixed(2)

						d3.select "#pie-tooltip"
							.style "left", "#{x}px"
							.style "top", "#{y}px"
							.attr "class", "tooltip fade top in"
							.select ".tooltip-inner"
								.text "#{d.data.key}: #{d.data.value} (#{percent}%)"

					slices.on "mouseout", (d) ->
						d3.select "#pie-tooltip"
							.attr "class", "tooltip top hidden"

		$scope.highlightPieSlice = (id, highlight) ->
			slice = d3.select "#slice-#{id}"
			return unless slice?
			slice.style 'opacity', if highlight is false then null else 0.75

		$scope.labelClass = (priority) ->
			return "label label-primary label-severe-#{priority + 1}"

		$scope.handleError = (type) ->
			(err, status) ->
				$scope.mainSections[type].error = switch
					when err and err.message then err.message
					when err and err.data and err.data.message then err.data.message
					when err and err.data then err.data
					when err then err
					else ''

				if $scope.mainSections[type].error.substring(0, 6) is '<html>'
					$scope.mainSections[type].error = do ->
						curr = $scope.mainSections[type].error
						curr = curr.replace(/(\r\n|\n|\r)/gm, '')
						curr = curr.replace(/.{0,}(\<title\>)/, '')
						curr = curr.replace(/(\<\/title\>).{0,}/, '')
						curr
		return
	]
