'use strict'

###*
 # @ngdoc function
 # @name dataplayApp.controller:SearchCtrl
 # @description
 # # SearchCtrl
 # Controller of the dataplayApp
###
angular.module('dataplayApp')
	.controller 'SearchCtrl', ['$scope', '$location', '$routeParams', 'User', 'Overview', 'PatternMatcher', ($scope, $location, $routeParams, User, Overview, PatternMatcher) ->
		$scope.query = if $routeParams.query? then $routeParams.query else ""
		$scope.searchTimeout = null
		$scope.results = []
		$scope.rowedResults = [] # split into sub-arrays of 3

		$scope.totalResults = 0
		$scope.rowLimit = 3
		$scope.overview = []

		# Charts
		$scope.allowed = ['line', 'bar', 'row', 'column', 'pie', 'bubble']
		$scope.params = $routeParams
		$scope.count = 3
		$scope.loading =
			related: false
		$scope.offset =
			related: 0
		$scope.limit =
			related: false
		$scope.max =
			related: 0
		$scope.chartsRelated = []

		$scope.xTicks = 6
		$scope.width = 275
		$scope.height = 200
		$scope.margin =
			top: 10
			right: 10
			bottom: 30
			left: 70

		$scope.init = (reset = false) ->
			# Initiate search if we have /search/:query
			if reset
				$scope.chartsRelated = []

			$scope.search()
			$scope.getNews()

		$scope.search = (offset = 0, count = 9) ->
			return if $scope.query.length < 3

			User.search $scope.query, offset, count
				.success (data) ->
					$scope.results = if offset is 0 then data.Results else $scope.results.concat data.Results

					$scope.results.forEach (r) ->
						r.graph = []
						r.error = null

						# Random
						offset = Overview.getRandomInteger 0, 3
						$scope.getRelated r.GUID, 0

					$scope.rowedResults = $scope.splitIntoRows $scope.results
					$scope.totalResults = data.Total
					return
				.error (status, data) ->
					console.log "Search::search::Error:", status
					return
			return

		$scope.getNews = () ->
			return if $scope.query.length < 3

			User.getNews $scope.query
				.success (data) ->
					if data instanceof Array
						$scope.overview = data.map (item) ->
							date: Overview.humanDate new Date item.date
							title: item.title
							url: item.url
							thumbnail: item['image_url']

				.error (status, data) ->
					console.log "Search::getNews::Error:", status
					return

		$scope.splitIntoRows = (arr, numOfCols = 3) ->
			twoD = []
			for item, key in arr
				row = Math.floor key / numOfCols
				col = key % numOfCols
				if not twoD[row]
					twoD[row] = []
				twoD[row][col] = item
			return twoD

		$scope.showMore = ->
			$scope.rowLimit += 2
			if $scope.rowedResults.length < $scope.rowLimit
				# get more results
				$scope.search ($scope.rowLimit - 1) * 3, 9

		$scope.collapse = (item) ->
			item.show = false

		$scope.uncollapse = (item) ->
			item.show = true


		# Charts
		$scope.findById = (id) ->
			data = _.where($scope.chartsRelated,
				id: id
			)

			if data?[0]? then data[0] else null

		$scope.isPlotAllowed = (type) ->
			if type in $scope.allowed then true else false

		$scope.getRelatedCharts = () ->
			$scope.getRelated Overview.charts 'related'

			return

		$scope.hasRelatedCharts = () ->
			Object.keys($scope.chartsRelated).length

		$scope.getRelated = (guid, offset) ->
			Overview.related guid, offset, 1
				.success (data) ->
					if data? and data.charts? and data.charts.length > 0
						for key, chart of data.charts
							continue unless $scope.isPlotAllowed chart.type

							key = parseInt(key)
							chart.key = key
							chart.id = "related-#{guid}-#{chart.key + $scope.offset.related}-#{chart.type}"
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

						console.log $scope.chartsRelated
					return
				.error (data, status) ->
					$scope.loading.related = false
					console.log "Overview::getRelated::Error:", status
					return

			return

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

		$scope.getXUnits = (data) ->
			xUnits = switch data.patterns[data.xLabel].valuePattern
				when 'date' then d3.time.years
				when 'intNumber' then dc.units.integers
				when 'label', 'text' then dc.units.ordinal
				else dc.units.ordinal

			xUnits

		$scope.getYScale = (data) ->
			yScale = switch data.patterns[data.yLabel].valuePattern
				when 'label'
					d3.scale.ordinal()
						.domain data.ordinals
						.rangeBands [0, $scope.height]
				when 'date'
					d3.time.scale()
						.domain d3.extent data.group.all(), (d) -> d.value
						.range [0, $scope.height]
				else
					d3.scale.linear()
						.domain d3.extent data.group.all(), (d) -> parseInt d.value
						.range [0, $scope.height]
						.nice()

			yScale

		$scope.lineChartPostSetup = (chart) ->
			data = $scope.findById chart.anchorName()

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> parseInt(d.y) % data.ordinals.length

			chart.xAxis().ticks $scope.xTicks

			chart.xAxisLabel false, 0
			chart.yAxisLabel false, 0

			chart.x $scope.getXScale data

			return

		$scope.rowChartPostSetup = (chart) ->
			data = $scope.findById chart.anchorName()

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.xAxis().ticks $scope.xTicks

			chart.x $scope.getYScale data

			chart.xUnits $scope.getXUnits data if data.ordinals?.length > 0

			return

		$scope.columnChartPostSetup = (chart) ->
			data = $scope.findById chart.anchorName()

			data.entry = crossfilter data.values
			data.dimension = data.entry.dimension (d) -> d.x
			data.group = data.dimension.group().reduceSum (d) -> d.y

			chart.dimension data.dimension
			chart.group data.group

			data.ordinals = []
			data.ordinals.push d.key for d in data.group.all() when d not in data.ordinals

			chart.colorAccessor (d, i) -> i + 1

			chart.xAxis().ticks $scope.xTicks

			chart.x $scope.getXScale data

			chart.xUnits $scope.getXUnits data if data.ordinals?.length > 0

			return

		$scope.pieChartPostSetup = (chart) ->
			data = $scope.findById chart.anchorName()

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

			return

		$scope.bubbleChartPostSetup = (chart) ->
			data = $scope.findById chart.anchorName()

			minR = null
			maxR = null

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
			for d in data.group.all() when d not in data.ordinals
				data.ordinals.push d.key.split("|")[0]

			chart.keyAccessor (d) -> d.key.split("|")[0]
			chart.valueAccessor (d) -> d.key.split("|")[1]
			chart.radiusValueAccessor (d) ->
				r = Math.abs parseInt d.key.split("|")[2]
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

			chart.xAxis().ticks $scope.xTicks

			# chart.label (d) -> x = d.key.split("|")[0]

			chart.title (d) ->
				x = d.key.split("|")[0]
				y = d.key.split("|")[1]
				z = d.key.split("|")[2]
				"#{data.xLabel}: #{x}\n#{data.yLabel}: #{y}\n#{data.zLabel}: #{z}"

			minRL = Math.log minR
			maxRL = Math.log maxR
			scale = Math.abs Math.log (maxRL - minRL) / (maxR - minR)

			chart.maxBubbleRelativeSize scale / 100

			return

		return
	]
