define ['jquery', 'underscore'], ($, _) ->
  'use strict'
  $ () ->    
    guid = window.location.href.split('/')[window.location.href.split('/').length - 1]
    template = if guid is "overlay" then '#searchResultsOverlayTemplate' else '#searchResultsTemplate'
    searchResultsTemplate = _.template $(template).html(), null, variable: 'data'
    $('#SBox').keyup () ->
      if $('#SBox').val()
        $.getJSON "/api/search/#{$('#SBox').val()}", (data) ->
          $('#ResultsTable').empty().append "<td><tr>Title</tr><tr>Guid</tr></td>"
          $('#ResultsTable').append(searchResultsTemplate data[i]) for i in [data.length-1..0]
