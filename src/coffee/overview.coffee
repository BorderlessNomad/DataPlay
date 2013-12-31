define ['jquery', 'app/PGPatternMatcher', 'app/PGOverviewCharts'], ($, PGPatternMatcher, PGOverviewCharts) ->
  'use strict'
  guid = window.location.href.split('/')[window.location.href.split('/').length - 1]
  $ () ->
    $.getJSON "/api/getdata/#{guid}", (data) ->
      if data.length
        patterns = {}
        for key of data[0]
          do (key) ->
            patterns[key] = PGPatternMatcher.getPattern data[0][key]
            # Now parse ALL the data
            entry[key] = PGPatternMatcher.parse(entry[key], patterns[key]) for entry in data
        new PGOverviewCharts guid, {dataset: data, patterns: patterns}, '#charts'
