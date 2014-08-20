'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCtrl
 # @description
 # # ChartsCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCtrl', ['$scope', '$routeParams', 'PatternMatcher', 'Charts', 'Tracker', ($scope, $routeParams, PatternMatcher, Charts, Tracker) ->
		$scope.params = $routeParams
		$scope.width = 1140;
		$scope.height = $scope.width * 9 / 16 # 16:9
		$scope.margin =
			top: 50
			right: 10
			bottom: 50
			left: 100
		$scope.marginAlt =
			top: 0
			right: 10
			bottom: 50
			left: 110

		$scope.chart =
			title: ""
			description: "N/A"
			data: null
			values: []
		$scope.monthNames = [
			"Jan"
			"Feb"
			"Mar"
			"Apr"
			"May"
			"Jun"
			"Jul"
			"Aug"
			"Sep"
			"Oct"
			"Nov"
			"Dec"
		]
		$scope.observations = []

		$scope.init = () ->
			# Track
			Tracker.visited $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z

			Charts.info $scope.params.id, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z
				.success (data) ->
					if data?
						$scope.chart = data

						if data.desc? and data.desc.length > 0
							description = data.desc.replace /(h1>|h2>|h3>)/ig, 'h4>'
							description = description.replace /\n/ig, ''
							$scope.chart.description = description

						$scope.reduceData()

					console.log "Chart", $scope.chart

					return
				.error (data, status) ->
					console.log "Charts::getInfo::Error:", status
					return

		$scope.reduceData = () ->
			$scope.chart.data = $scope.chart.values
			$scope.chart.patterns = {}
			$scope.chart.patterns[$scope.chart.xLabel] =
				valuePattern: PatternMatcher.getPattern $scope.chart.values[0]['x']
				keyPattern: PatternMatcher.getKeyPattern $scope.chart.values[0]['x']

			if $scope.chart.patterns[$scope.chart.xLabel].valuePattern is 'date'
				for value, key in $scope.chart.values
					$scope.chart.values[key].x = new Date(value.x)

			if $scope.chart.yLabel?
				$scope.chart.patterns[$scope.chart.yLabel] =
					valuePattern: PatternMatcher.getPattern $scope.chart.values[0]['y']
					keyPattern: PatternMatcher.getKeyPattern $scope.chart.values[0]['y']

		$scope.humanDate = (date) ->
			"#{date.getDate()} #{$scope.monthNames[date.getMonth()]}, #{date.getFullYear()}"

		$scope.getXScale = (data) ->
			xScale = switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.width]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key
						.range [0, $scope.width]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key
						.range [0, $scope.width]

			xScale

		$scope.sortY = (a, b) ->
			if $scope.chart.patterns[$scope.chart.yLabel].valuePattern is 'date'
				a.value = a.value.getTime()
				b.value = b.value.getTime()

			a.value - b.value

		$scope.getYScale = (data, group) ->
			yScale = switch data.patterns[data.yLabel].valuePattern
				when 'label'
					unless group?
						group = data.ordinals

					d3.scale.ordinal()
						.domain group
						.rangeBands [$scope.height, 0]
				when 'date'
					unless group?
						group = data.group.all()

					if sort? and sort is true
						group = Array.prototype.slice.call(group).sort(compare)

					d3.time.scale()
						.domain d3.extent group, (d) -> d.value
						.range [$scope.height, 0]
				else
					unless group?
						group = data.group.all()

					d3.scale.linear()
						.domain d3.extent group, (d) -> parseInt d.value
						.range [$scope.height, 0]

			yScale

		$scope.nearestNeighbour = (index, data, key) ->
			index = parseInt index
			curr = data[index][key]
			next = data[index + 1][key]
			prev = data[index - 1][key]

			closest = curr
			if curr is next
				closest = next
			else if curr is prev
				closest = prev
			else if Math.abs(curr - next) >= Math.abs(curr - prev)
				closest = prev
			else if Math.abs(curr - next) <= Math.abs(curr - prev)
				closest = next

			closest

		$scope.drawCircle = (area, data, x, y, plot, color) ->
			area.append 'circle'
				.attr 'cx', plot[0]
				.attr 'cy', plot[1]
				.attr 'r', 5
				.style 'fill', color
				.append 'svg:title'
				.text (d) ->
					if data.patterns[data.xLabel].valuePattern is 'date'
						x = $scope.humanDate x

					"#{$scope.chart.xLabel}: #{x}\n#{$scope.chart.yLabel}: #{y}"

		$scope.lineChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y
			# data.group1 = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group, data.title
			# chart.stack data.group1, "Test"

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			# chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length
			chart.keyAccessor (d) -> d.key
			chart.valueAccessor (d) -> d.value
			chart.title (d) ->
				x = d.key
				if $scope.chart.patterns[$scope.chart.xLabel].valuePattern is 'date'
					x = $scope.humanDate d.key
				"#{$scope.chart.xLabel}: #{x}\n#{$scope.chart.yLabel}: #{d.value}"
			chart.legend dc.legend().itemHeight(13).gap(5)

			xScale = $scope.getXScale data
			chart.x xScale

			yScaleGroup = Array.prototype.slice.call(data.group.all()).sort($scope.sortY)
			yScale = $scope.getYScale data

			if data.ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			points = [
				[1950, 600000]
				[1960, 800000]
				[1970, 700000]
			]

			chart.renderlet (c) ->
				console.log 'renderlet'

				svg = d3.select 'svg'
				stack = d3.select('g.stack-list').node()
				box = stack.getBBox()

				area = d3.select 'g.stack-list'
					.append 'g'
					.attr 'class', 'stack _1' # change

				circles = c.svg().selectAll 'circle.dot'
				circleTitles = c.svg().selectAll 'circle.dot > title'
				yDomain = []
				for cr, i in circles[0]
					yDomain[i] =
						cy: cr.cy.animVal.value
						y: null

				for crt, i in circleTitles[0]
					html = crt.innerHTML.split('\n')
					yDomain[i].y = html[1].split(': ')[1]

				yDomain.sort (a, b) -> a.cy - b.cy

				# Exisiting observations
				for p in points
					x = p[0]
					y = p[1]
					plot = [xScale(x), yScale(y)]
					color = '#ff7f0e'
					console.log x, y, plot, color

					$scope.drawCircle area, data, x, y, plot, color

				# New observations
				datum = null
				circles.on 'click', (d) ->
					console.log 'POINT', 'Datum:', d
					datum = d

				svg.on 'click', () ->
					space = d3.mouse(stack)

					# Clicking out-side box is now allowed
					if space[0] < 0 or space[1] < 0 or space[0] > box.width or space[1] > box.height
						return

					color = '#2ca02c'
					if datum?
						color = '#d62728'
						x = datum.x
						y = datum.y
						datum = null
					else
						x = xScale.invert space[0]
						y = yScale.invert space[1]

					# X
					for val, key in data.group.all()
						if data.patterns[data.xLabel].valuePattern is 'date'
							if x.getTime() >= val.key.getTime()
								i = key
						else if x >= val.key
							i = key

					x = data.group.all()[i].key
					$scope.nearestNeighbour i, data.group.all(), 'key'

					# Y
					for k, v of yDomain
						if v.cy <= space[1]
							j = k

					y = yDomain[j].y
					$scope.nearestNeighbour j, yDomain, 'y'

					$scope.observations.push
						x: x
						y: y
						space: space
						color: color

					$scope.drawCircle area, data, x, y, space, color

			return

		$scope.rangeChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.x $scope.getXScale data

			if data.ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			return

		$scope.rowChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> (i + 1) % 20

			chart.xAxis()
				.ticks $scope.xTicks

			chart.x $scope.getYScale data

			if ordinals? and ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			return

		$scope.columnChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.x $scope.getXScale data

			if ordinals? and ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal

			return

		$scope.pieChartPostSetup = (chart) ->
			data = $scope.chart

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) ->
				if data.patterns[data.xLabel].valuePattern is 'date'
					return "#{d.x.getDate()} #{$scope.monthNames[d.x.getMonth()]} #{d.x.getFullYear()}"
				x = if d.x? and (d.x.length > 0 || data.patterns[data.xLabel].valuePattern is 'date') then d.x else "N/A"
			data.groupSum = 0
			data.group = data.dimension.group().reduceSum (d) ->
				y = Math.abs parseFloat d.y
				data.groupSum += y
				y

			chart.dimension data.dimension
			chart.group data.group

			chart.colorAccessor (d, i) -> i + 1

			chart.renderLabel false
			chart.label (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key} (#{Math.floor percent}%)"

			chart.renderTitle false
			chart.title (d) ->
				percent = d.value / data.groupSum * 100
				"#{d.key}: #{d.value} [#{Math.floor percent}%]"

			chart.legend dc.legend()

			chart.minAngleForLabel 0
			chart.innerRadius 150

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.chart

			minR = null
			maxR = null

			###
			# TODO: when X/Y is String group all similar X/Y and SUM Y/X.
			# if Z is 'transaction_number' or 'invoice_number' ignore it
			# and replace radius with count of X/Y
			###

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) ->
				z = Math.abs parseInt d.z

				if not minR? or minR > z
					minR = if z is 0 then 1 else z

				if not maxR? or maxR <= z
					maxR = if z is 0 then 1 else z

				"#{d.x}|#{d.y}|#{d.z}"

			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key.split("|")[0] for d in data.group.all() when d not in data.ordinals
			console.log data.ordinals

			chart.keyAccessor (d) -> d.key.split("|")[0]
			chart.valueAccessor (d) -> d.key.split("|")[1]
			chart.radiusValueAccessor (d) ->
				r = Math.abs d.key.split("|")[2]
				if r >= minR then r else minR

			chart.x switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.width]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key.split("|")[0]
						.range [0, $scope.width]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[0]
						.range [0, $scope.width]

			chart.y switch data.patterns[data.xLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.height]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.key.split("|")[1]
						.range [0, $scope.height]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.key.split("|")[1]
						.range [0, $scope.height]

			rScale = d3.scale.linear()
				.domain d3.extent data.group.all(), (d) -> Math.abs parseInt d.key.split("|")[2]
			chart.r rScale

			chart.label (d) -> x = d.key.split("|")[0]

			chart.title (d) ->
				x = d.key.split("|")[0]
				y = d.key.split("|")[1]
				z = d.key.split("|")[2]
				"#{data.xLabel}: #{x}\n#{data.yLabel}: #{y}\n#{data.zLabel}: #{z}"

			minRL = Math.log minR
			maxRL = Math.log maxR
			scale = Math.abs Math.log (maxRL - minRL) / (maxR - minR)

			console.log minR, maxR, (maxR - minR), minRL, maxRL, (maxRL - minRL), scale, scale / 100

			chart.maxBubbleRelativeSize scale / 100

			chart.zoomOutRestrict true
			chart.mouseZoomable true

			# chart.renderlet (c) ->
			# 	circles = c.svg().selectAll('g.chart-body').selectAll('g circle')
			# 	for circle in circles[0]
			# 		r = circle.r
			# 		if r.baseVal.value < 0
			# 			r.baseVal.value = Math.abs r.baseVal.value
			# 			r.animVal = r.baseVal

			return

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

		return
	]
