define ['jquery', 'underscore'], ($, _) ->
  'use strict'
  $ () ->
    visitedEntryTemplate = _.template $('#visitedEntryTemplate').html(), null, variable: 'data'
    $.getJSON "/api/user", (data) ->
      $('#FillInUserName').text data.Username
      $.getJSON "/api/visited", (data) ->
        if data and data.length
          $('#History').empty().append "<p>Welcome Back, You where last viewing:</p>"
          $('#History').append(visitedEntryTemplate entry) for entry in data
