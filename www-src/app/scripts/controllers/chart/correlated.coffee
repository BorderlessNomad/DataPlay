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
		$scope.mode = 'correlated'
		$scope.width = 570
		$scope.height = $scope.width * 9 / 16 # 16:9
		$scope.margin =
			top: 50
			right: 20
			bottom: 50
			left: 100
		$scope.marginAlt =
			top: 0
			right: 20
			bottom: 50
			left: 110
		$scope.xTicks = 8

		$scope.chartRendered = null
		$scope.chart =
			title: ""
			description: "N/A"
			data: null
			values: []
		$scope.xScale = null
		$scope.yDomain = null
		$scope.userObservations = []
		$scope.userObservationsMessage = []
		$scope.observation =
			x: null
			y: null
			message: ''

		$scope.info =
			discoveredId: null
			validated: null
			invalidated: null
			patternId: null
			discoverer: ''
			discoverDate: ''
			validators: []
			source:
				prim: ''
				seco: ''
			strength: ''

		$scope.init = () ->
			$scope.initChart()
			return

		$scope.initChart = () ->
			Charts.correlated $scope.params.correlationid
				.success (data, status) ->
					if data? and data.chartdata
						[1..2].forEach (i) ->
							vals = $scope.translateToNv data.chartdata['table' + i].values
							dataRange = do ->
								min = d3.min vals, (item) -> parseFloat item.y
								[
									if min > 0 then 0 else min
									d3.max vals, (item) -> parseFloat item.y
								]
							$scope.corrLine.data.push
								key: data.chartdata['table' + i].title
								type: 'area'
								yAxis: i
								values: vals
							$scope.corrLine.options.chart['yDomain' + i] = dataRange
							$scope.corrLine.options.chart['yAxis' + i].tickValues = do ->
								[1..8].map (num) ->
									dataRange[0] + ((dataRange[1] - dataRange[0]) * ((1 / 8) * num))

					if data?
						$scope.info.patternId = data.patternid or ''
						$scope.info.discoveredId = data.discoveredid or ''
						$scope.info.discoverer = data.discoveredby or ''
						$scope.info.discoverDate = if data.discoverydate then Overview.humanDate new Date( data.discoverydate ) else ''
						$scope.info.validators = data.validatedby or ''
						$scope.info.source =
							prim: data.source1 or ''
							seco: data.source2 or ''
						$scope.info.strength = data.statstrength

					$scope.initObservations()
					console.log "Chart", $scope.chart

					# Track a page visit
					Tracker.visited $scope.params.id, $scope.params.key, $scope.params.type, $scope.params.x, $scope.params.y, $scope.params.z
				.error (data, status) ->
					console.log "Charts::init::Error:", status

			return

		$scope.initObservations = (redraw) ->
			id = "#{$scope.params.correlationid}"

			Charts.getObservations $scope.info.discoveredId
				.then (res) ->
					$scope.userObservations.splice 0, $scope.userObservations.length

					res.data?.forEach? (obsv) ->
						x = obsv.x
						y = obsv.y
						if $scope.chart.patterns[$scope.chart.table1.xLabel].valuePattern is 'date'
							if not(x instanceof Date) and (typeof x is 'string')
								xdate = new Date x
								if xdate.toString() isnt 'Invalid Date' then x = xdate
							x = Overview.humanDate x

						xy = "#{x.replace(/\W/g, '')}-#{y.replace(/\W/g, '')}"
						$scope.userObservationsMessage[xy] = obsv.comment
						$scope.userObservations.push
							xy: xy
							oid : obsv['observation_id']
							user: obsv.user
							validationCount: parseInt(obsv.validations - obsv.invalidations) || 0
							message: obsv.comment
							date: Overview.humanDate new Date(obsv.created)
							coor:
								x: obsv.x
								y: obsv.y

				, $scope.handleError
			return

		$scope.translateToNv = (values) ->
			normalise = (d) ->
				if typeof d is 'string'
					if not isNaN Date.parse d
						return Date.parse d
					if not isNaN parseFloat d
						return parseFloat d
				return d

			values.map (v) ->
				x: normalise v.x || 0
				y: parseFloat v.y || 0

		$scope.corrLine =
			options:
				chart:
					type: "multiChart"
					height: 450
					margin:
						top: 15
						right: 35
						bottom: 25
						left: 35
					x: (d, i) -> i
					y: (d) -> d[1]
					color: d3.scale.category10().range()
					transitionDuration: 250
					xAxis:
						axisLabel: ""
						showMaxMin: false
						tickFormat: (d) -> d3.time.format("%d-%m-%Y") new Date d
						ticks: $scope.xTicks
					yAxis1:
						orient: 'left'
						axisLabel: ""
						tickFormat: (d) -> d3.format(",f") d
						showMaxMin: false
						highlightZero: false
					yAxis2:
						orient: 'right'
						axisLabel: ""
						tickFormat: (d) -> d3.format(",f") d
						showMaxMin: false
						highlightZero: false
					areas:
						dispatch:
							elementClick: (e) ->
								console.log 'areas elementClick', e
					lines:
						dispatch:
							elementClick: (e) ->
								console.log 'lines elementClick', e
					yDomain1: [0, 1000]
					yDomain2: [0, 1000]
			data: []

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

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

			$scope.resetObservations()

		return
	]
