'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:ChartsCtrl
 # @description
 # # ChartsCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'ChartsCtrl', ['$scope', '$routeParams', 'Charts', 'PatternMatcher', ($scope, $routeParams, Charts, PatternMatcher) ->
		$scope.params = $routeParams
		$scope.info = {}
		$scope.chart = {}
		$scope.cfdata = null

		$scope.init = () ->
			$scope.getInfo()
			$scope.getData()
			return

		$scope.getInfo = () ->
			Charts.info $scope.params.id
				.success (data) ->
					if data?
						$scope.info = data
					return
				.error (data) ->
					console.log "Charts::getInfo::Error:", status
					return

		$scope.getData = () ->
			Charts.reducedData $scope.params.id, $scope.params.x, $scope.params.y, 10, 100
				.success (data) ->
					if data? and data.length
						$scope.chart.dataset = data
						$scope.chart.keys = []
						$scope.chart.patterns = {}

						for key of data[0]
							do (key) ->
								$scope.chart.keys.push key
								$scope.chart.patterns[key] =
									valuePattern: PatternMatcher.getPattern data[0][key]
									keyPattern: PatternMatcher.getKeyPattern data[0][key]

						console.log "getData", $scope.chart

						$scope.identifyData()
					return
				.error (data) ->
					console.log "Charts::getChart::Error:", status
					return

		$scope.identifyData = () ->
			Charts.identifyData $scope.params.id
				.success (data) ->
					for col in data.Cols
						do (col) ->
							switch col.Sqltype
								when "int", "bigint"
									$scope.chart.patterns[col.Name].valuePattern = 'intNumber' if $scope.chart.patterns?[col.Name]?.valuePattern?
								when "float"
									$scope.chart.patterns[col.Name].valuePattern = 'floatNumber' if $scope.chart.patterns?[col.Name]?.valuePattern?

					console.log "identifyData", $scope.chart

					$scope.initChart 'lines'
					return
				.error (data) ->
					console.log "Charts::identifyData::Error:", status
					return

		$scope.initChart = (type) ->
			if $scope.chart.type isnt type
				$scope.chart.type = type
				$scope.chart.data = null
				# $scope.populateKeys()
				$scope.parseChartData $scope.params.id, $scope.params.x, $scope.params.y, $scope.chart.type, $scope.chart.dataset

		$scope.populateKeys = () ->
			# TODO: Dropdown of X & Y axis

		$scope.parseChartData = (guid, x, y, type, data) ->
			if type in ['bars', 'pie', 'bubbles']
				Charts.groupedData guid, x, y
					.success (results) ->
						$scope.chart.data = $scope.parseResults x, y, results
						$scope.cfdata = crossfilter $scope.chart.dataset
						console.log "parseChartData", $scope.chart.data
						return
					.error (results) ->
						console.log "Charts::parseChartData::Error:", status
						return
			else
				$scope.chart.data = $scope.parseResults x, y, data
				console.log "parseChartData", $scope.chart.data

			return

		$scope.parseResults = (xAxis, yAxis, results) ->
			datapool = []

			console.log "parseResults", xAxis, yAxis, results, $scope.chart.patterns[xAxis], $scope.chart.patterns[yAxis]

			x =
				key: xAxis
				pattern: $scope.chart.patterns[xAxis]
			y =
				key: yAxis
				pattern: $scope.chart.patterns[yAxis]

			xData = $scope.parseAxisData results, x
			yData = $scope.parseAxisData results, y

			datapool.push [xData[i] , yData[i]] for i in [0..results.length - 1]

			datapool

		$scope.parseAxisData = (results, axis) ->
			data = []
			valuePattern = PatternMatcher.getPattern results[0][axis.key]
			for item in results
				do (item) =>
					data.push PatternMatcher.parse item[axis.key], valuePattern

			data

		$scope.lineChartPostSetup = (chart, options) ->
			console.log "lineChartPostSetup", chart, options

			chart.colorAccessor (d) -> parseInt(d.value) % 20

			chart.x

			return

		return
	]
