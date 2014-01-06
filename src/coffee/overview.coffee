define ['jquery', 'app/PGPatternMatcher', 'app/PGOverviewCharts'], ($, PGPatternMatcher, PGOverviewCharts) ->
  'use strict'
  guid = window.location.href.split('/')[window.location.href.split('/').length - 1]
  $ () ->
    $.getJSON "/api/getreduceddata/#{guid}/15/100", (data) ->
      if data.length
        patterns = {}
        for key of data[0]
          do (key) ->
            vp = PGPatternMatcher.getPattern data[0][key]
            kp = PGPatternMatcher.getKeyPattern key
            patterns[key] = valuePattern: vp, keyPattern: kp
            # Now parse ALL the data
            # TODO get into account key pattern before parsing everything???
            entry[key] = PGPatternMatcher.parse(entry[key], patterns[key].valuePattern) for entry in data
        new PGOverviewCharts guid, {dataset: data, patterns: patterns}, '#charts'
