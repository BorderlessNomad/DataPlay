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

		$scope.corrChart = new CorrelatedChart

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
						$scope.corrChart.generate data.chartdata.type

						if not $scope.corrChart.error
							[1..2].forEach (i) ->
								vals = $scope.translateToNv data.chartdata['table' + i].values, data.chartdata.type
								dataRange = do ->
									min = d3.min vals, (item) -> parseFloat item.y
									[
										if min > 0 then 0 else min
										d3.max vals, (item) -> parseFloat item.y
									]

								$scope.corrChart.data.push
									key: data.chartdata['table' + i].title
									type: 'area'
									yAxis: i
									values: vals
								$scope.corrChart.options.chart['yDomain' + i] = dataRange
								$scope.corrChart.options.chart['yAxis' + i].tickValues = do ->
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

		$scope.translateToNv = (values, type) ->
			normalise = (d) ->
				if typeof d is 'string'
					if not isNaN Date.parse d
						return Date.parse d
					if not isNaN parseFloat d
						return parseFloat d
				return d

			values.map (v) ->
				newV =
					x: normalise v.x || 0
					y: parseFloat v.y || 0
				if type is 'scatter'
					newV.size = 0.5
					newV.shape = 'circle'
				newV

		$scope.validateChart = (valFlag) ->
			id = "#{$scope.params.id}/#{$scope.params.key}/#{$scope.params.type}/#{$scope.params.x}/#{$scope.params.y}"
			id += "/#{$scope.params.z}" if $scope.params.z?.length > 0

			Charts.validateChart "rid", id, valFlag
				.then ->
					$scope.showValidationMessage valFlag
					if valFlag
						$scope.info.validated = true
					else
						$scope.info.invalidated = true
				, $scope.handleError

		$scope.saveObservation = ->
			Charts.createObservation($scope.info.discoveredId, $scope.observation.x, $scope.observation.y, $scope.observation.message).then (res) ->
				$scope.observation.message = ''

				$scope.addObservation $scope.observation.x, $scope.observation.y, $scope.observation.message

				$('#comment-modal').modal 'hide'
			, $scope.handleError

			return

		$scope.clearObservation = ->
			$scope.observation.message = ''

			x = $scope.observation.x
			y = $scope.observation.y

			if (x is 0 or x is "0") and (y is 0 or y is "0")
				$('#comment-modal').modal 'hide'
				return

			if not(x instanceof Date) and (typeof x is 'string')
				xdate = new Date x
				if xdate.toString() isnt 'Invalid Date' then x = Overview.humanDate xdate
			else if x instanceof Date
				x = Overview.humanDate x

			xy = "#{x.replace(/\W/g, '')}-#{y.replace(/\W/g, '')}"

			# console.log xy, d3.select("#clipImage-#{xy}"), d3.select("#observationIcon-#{xy}")
			d3.select("#clipImage-#{xy}").remove()
			d3.select("#observationIcon-#{xy}").remove()

			$('#comment-modal').modal 'hide'
			return

		$scope.validateObservation = (item, valFlag) ->
			if item.oid?
				Charts.validateObservation item.oid, valFlag
					.success (res) ->
						item.validationCount += (valFlag) ? 1 : -1
					.error $scope.handleError

		$scope.openAddObservationModal = (x, y) ->
			$scope.observation.x = x || 0
			$scope.observation.y = y || 0
			$scope.observation.message = ''

			$('#comment-modal').modal 'show'

			return

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
