'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCorrelatedCtrl
 # @description
 # # ChartsCorrelatedCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCorrelatedCtrl', ['$scope', '$location', '$timeout', '$routeParams', 'Auth', 'config', 'Overview', 'PatternMatcher', 'Charts', ($scope, $location, $timeout, $routeParams, Auth, config, Overview, PatternMatcher, Charts) ->
		$scope.username = Auth.get config.userName

		$scope.params = $routeParams
		$scope.mode = 'correlated'
		$scope.width = 570
		$scope.height = $scope.width * 9 / 16 # 16:9

		$scope.chart = new CorrelatedChart

		$scope.userObservations = null
		$scope.userObservationsMessage = []
		$scope.observation =
			x: null
			y: null
			message: ''

		$scope.info =
			discoveredId: null
			credited: false
			discredited: false
			patternId: null
			discoverer: ''
			discoverDate: ''
			creditors: []
			discreditors: []
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
						$scope.chart.generate data.chartdata.type

						if not $scope.chart.error
							[1..2].forEach (i) ->
								vals = $scope.chart.translateData data.chartdata['table' + i].values, data.chartdata.type
								dataRange = do ->
									min = d3.min vals, (item) -> parseFloat item.y
									[
										if min > 0 then 0 else min
										d3.max vals, (item) -> parseFloat item.y
									]
								type = if data.chartdata.type is 'column' or data.chartdata.type is 'bar' then 'bar' else 'area'

								$scope.chart.data.push
									key: data.chartdata['table' + i].title
									type: type
									yAxis: i
									values: vals
								$scope.chart.options.chart['yDomain' + i] = dataRange
								$scope.chart.options.chart['yAxis' + i].tickValues = do ->
									[1..8].map (num) ->
										dataRange[0] + ((dataRange[1] - dataRange[0]) * ((1 / 8) * num))

							$scope.chart.setAxisTypes data.chartdata.table1.xLabel, data.chartdata.table1.yLabel, data.chartdata.table2.yLabel

					if data?
						$scope.info.patternId = data.patternid or ''
						$scope.info.discoveredId = data.discoveredid or ''
						$scope.info.discoverer = data.discoveredby or ''
						$scope.info.discoverDate = if data.discoverydate then Overview.humanDate new Date( data.discoverydate ) else ''
						$scope.info.creditors = data.creditedby or ''
						$scope.info.discreditors = data.discreditedby or ''
						$scope.info.source =
							prim: data.source1 or ''
							seco: data.source2 or ''
						$scope.info.strength = data.statstrength
						$scope.info.credited = data.userhascredited
						$scope.info.discredited = data.userhasdiscredited

					$scope.initObservations()
					console.log "Chart", $scope.chart
				.error (data, status) ->
					console.log "Charts::init::Error:", status

			return

		$scope.initObservations = (redraw) ->
			Charts.getObservations $scope.info.discoveredId
				.then (res) ->
					$scope.userObservations = []

					res.data?.forEach? (obsv) ->
						x = "0"
						y = "0"
						xy = "#{x.replace(/\W/g, '')}-#{y.replace(/\W/g, '')}"
						$scope.userObservationsMessage[xy] = obsv.comment
						if obsv.user.avatar is ''
							obsv.user.avatar = "http://www.gravatar.com/avatar/#{obsv.user.email}?d=identicon"
						$scope.userObservations.push
							xy: xy
							oid : obsv['observation_id']
							user: obsv.user
							credits: obsv.credits
							discredits: obsv.discredits
							creditCount: parseInt(obsv.credits - obsv.discredits) || 0
							message: obsv.comment
							date: Overview.humanDate new Date(obsv.created)
							coor:
								x: obsv.x
								y: obsv.y
							flagged: !! obsv.flagged
							action: obsv.action

				, $scope.handleError
			return

		$scope.creditChart = (valFlag) ->
			Charts.creditChart "cid", $scope.params.correlationid, valFlag
				.then ->
					$scope.showCreditMessage valFlag
					$scope.info.credited = !! valFlag
					$scope.info.discredited = ! valFlag
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

		$scope.creditObservation = (item, valFlag) ->
			if item.oid?
				Charts.creditObservation item.oid, valFlag
					.success (res) ->
						item.credits = res.Credited
						item.discredits = res.Discredited
						item.creditCount = parseInt(res.credits - res.discredits) || 0
						item.action = res.action
						item.flagged = !! res.flagged
					.error $scope.handleError

		$scope.openAddObservationModal = (x, y) ->
			$scope.observation.x = x || 0
			$scope.observation.y = y || 0
			$scope.observation.message = ''

			$('#comment-modal').modal 'show'

			return

		$scope.addObservation = (x, y, comment) ->
			$scope.initObservations(true)

		$scope.resetObservations = ->
			$scope.observations = []

		$scope.flagObservation = (obsv) ->
			Charts.flagObservation obsv.oid
				.success (data) ->
					obsv.flagged = true

		$scope.resetAll = ->
			dc.filterAll()
			dc.redrawAll()

			$scope.resetObservations()

		$scope.showCreditMessage = (type) ->
			$scope.creditMsg = type
			$timeout ->
				$scope.creditMsg = null
			, 3000
			return

		$scope.handleError = (err, status) ->
			$scope.error =
				message: switch
					when err and err.message then err.message
					when err and err.data and err.data.message then err.data.message
					when err and err.data then err.data
					when err then err
					else ''

			if $scope.error.substring(0, 6) is '<html>'
				$scope.error = do ->
					curr = $scope.error
					curr = curr.replace(/(\r\n|\n|\r)/gm, '')
					curr = curr.replace(/.{0,}(\<title\>)/, '')
					curr = curr.replace(/(\<\/title\>).{0,}/, '')
					curr

		return
	]
