'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCorrelatedCtrl
 # @description
 # # ChartsCorrelatedCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCorrelatedCtrl', ['$scope', '$location', '$routeParams', 'Overview', 'PatternMatcher', 'Charts', 'Tracker', ($scope, $location, $routeParams, Overview, PatternMatcher, Charts, Tracker) ->
		$scope.params = $routeParams
		$scope.width = 570
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
		$scope.observations = []
		$scope.userObservations = []
		$scope.observation =
			x: null
			y: null
			message: ''

		$scope.info =
			patternId: '202121200'
			discoverer: 'DataWiz'
			discoverDate: Overview.humanDate new Date( new Date() - (2 * 24 * 60 * 60 * 1000) )
			validators: [
				'Alan'
				'Bob'
				'Chris'
			]
			source:
				prim: 'NHS Spending 2012 - London'
				seco: 'Weather Patterns 2012'
			strength: 'High'

		$scope.init = () ->
			Charts.correlated $scope.params.correlationid
				.success (data, status) ->
					if data?
						$scope.chart = data

						if data.desc? and data.desc.length > 0
							description = data.desc.replace /(h1>|h2>|h3>)/ig, 'h4>'
							description = description.replace /\n/ig, ''
							$scope.chart.description = description

						$scope.reduceData()

					console.log "Chart", $scope.chart

					# Track a page visit
					Tracker.visited $scope.params.id, $scope.params.key, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z
				.error (data, status) ->
					console.log "Charts::init::Error:", status

			Charts.validateChart "#{$scope.params.correlationid}"
				.then (validate) ->
					valId = validate.data
					Charts.getObservations valId
						.then (res) ->
							$scope.userObservations.splice 0, $scope.userObservations.length

							res.data?.forEach (obsv) ->
								$scope.userObservations.push
									oid : obsv['observation_id']
									user: obsv.user
									validationCount: parseInt(obsv.validations - obsv.invalidations) || 0
									message: obsv.comment
									date: Overview.humanDate new Date(obsv.created)
									coor:
										x: obsv.x
										y: obsv.y

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

		$scope.getYScale = (data) ->
			yScale = switch data.patterns[data.yLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [$scope.height, 0]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.value
						.range [$scope.height, 0]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.value
						.range [$scope.height, 0]

			yScale

		$scope.nearestNeighbour = (index, data, key) ->
			index = parseInt index
			curr = data[index][key]
			next = if data[index + 1]? then data[index + 1][key] else null
			prev = if data[index - 1]? then data[index - 1][key] else null

			closest = curr
			if curr is next
				closest = next
			else if curr is prev
				closest = prev
			else if not next? or Math.abs(curr - next) >= Math.abs(curr - prev)
				closest = prev
			else if not prev? or Math.abs(curr - next) <= Math.abs(curr - prev)
				closest = next

			closest

		$scope.drawCircle = (area, data, x, y, plot, color) ->
			if data.patterns[data.xLabel].valuePattern is 'date'
				x = Overview.humanDate x

			circ = area.append 'circle'
				.attr 'cx', plot[0]
				.attr 'cy', plot[1]
				.attr 'r', 5
				.style 'fill', color
				.attr 'data-placement', 'top'
				.attr 'data-html', true
				.tooltip "#{$scope.chart.xLabel}: #{x}<br/>#{$scope.chart.yLabel}: #{y}"

		$scope.saveObservation = ->
			# console.log Charts.createObservation($scope.params.type, $scope.observation.x, $scope.observation.y, $scope.observation.message).then (res) ->
			# 	console.log 'suc', res
			# , (err) ->
			# 	console.log 'err', err

			$scope.observation.message = ''
			$('#comment-modal').modal 'hide'
			return

		$scope.clearObservation = ->
			$scope.observation.message = ''
			$('#comment-modal').modal 'hide'
			return


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
					x = Overview.humanDate d.key
				"#{$scope.chart.xLabel}: #{x}\n#{$scope.chart.yLabel}: #{d.value}"
			# chart.legend dc.legend().itemHeight(13).gap(5)

			xScale = $scope.getXScale data
			chart.x xScale

			yScale = $scope.getYScale data

			if data.ordinals.length > 0
				chart.xUnits switch data.patterns[data.xLabel].valuePattern
					when 'date' then d3.time.years
					when 'intNumber' then dc.units.integers
					when 'label', 'text' then dc.units.ordinal
					else dc.units.ordinal


			points = [
				[new Date("Feb 01 1975 00:00:00 GMT+0000 (GMT Standard Time)"), 600, "Hello Jack!"]
				[new Date("Jan 04 2000 00:00:00 GMT+0000 (GMT Standard Time)"), 400, "Test Point"]
			]

			existingObservations = null
			newObservations = null

			chart.renderlet (c) ->
				console.log 'renderlet'

				svg = d3.select 'svg'
				stack = d3.select('g.stack-list').node()
				box = stack.getBBox()

				stackList = d3.select 'g.stack-list'
				if not existingObservations?
					existingObservations = stackList.append 'g'
						.attr 'class', "stack _#{stackList.length + 0} observations existing"

				if not newObservations?
					newObservations = stackList.append 'g'
						.attr 'class', "stack _#{stackList.length + 1} observations new"

				# Clean observation points and re-render
				d3.selectAll('g.observations > circle').remove()

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

					# Y
					for k, v of yDomain
						if v.y is p[1]
							j = k
							break
						else if parseInt(v.y) > parseInt(p[1])
							j = k

					# Do not consider points which are beyond visible Y-axis
					if not yDomain? or not yDomain[j]? or not yDomain[j].y?
						continue

					y = yDomain[j].y

					plot = [xScale(x), yDomain[j].cy]
					color = '#ff7f0e'

					$scope.drawCircle existingObservations, data, x, y, plot, color

				# New observations
				datum = null
				circles.on 'click', (d) ->
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
							if x.getTime() is val.key.getTime()
								i = key
								break
							else if x.getTime() > val.key.getTime()
								i = key
						else if x is val.key
							i = key
							break
						else if x > val.key
							i = key

					x = data.group.all()[i].key
					$scope.nearestNeighbour i, data.group.all(), 'key'

					# Y
					for k, v of yDomain
						if v.cy is space[1]
							j = k
							break
						else if v.cy <= space[1]
							j = k

					y = yDomain[j].y
					$scope.nearestNeighbour j, yDomain, 'y'

					$scope.addObservation x, y, space

					$scope.observation.x = x
					$scope.observation.y = y
					$scope.observation.message = ''

					$('#comment-modal').modal 'show'
					$scope.drawCircle newObservations, data, x, y, space, color

				for p in $scope.observations
					x = p.x

					# Y
					for k, v of yDomain
						if v.y is p.y
							j = k
							break
						else if parseInt(v.y) > parseInt(p.y)
							j = k

					# Do not consider points which are beyond visible Y-axis
					if not yDomain? or not yDomain[j]? or not yDomain[j].y?
						continue

					y = yDomain[j].y

					plot = [xScale(x), yDomain[j].cy]
					color = '#2ca02c'

					$scope.drawCircle newObservations, data, x, y, plot, color

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
					return Overview.humanDate d.x
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
			chart.innerRadius 75

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

		$scope.validateObservation = (item, valFlag) ->
			if item.oid?
				Charts.validateObservation item.oid, valFlag
					.success (res) ->
						item.validationCount += (valFlag) ? 1 : -1

		$scope.addObservation = (x, y, space, comment) ->
			$scope.observations.push
				x: x
				y: y
				space: space
				comment: if comment? and comment.length > 0 then comment else ""
				timestamp: Date.now()

			$scope.$apply()

		$scope.resetObservations = ->
			$scope.observations = []

			d3.selectAll('g.observations.new > circle').remove()

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

			$scope.resetObservations()

		return
	]
