define ['jquery', 'typeahead', 'mustache'], ($, Typeahead, Mustache) ->
  'use strict'
  guids = []
  $ () ->    
    Typeahead(
      $('#dc-search .dc-search-input input')
      source: (query, result) ->
        $.getJSON "/api/search/#{$(@).val()}", (data) -> 
          guids[data[i].Title] = data[i].GUID for i in [data.length-1..0]
          titles = (data[i].Title for i in [data.length-1..0])
          result titles
      updater: (item) -> location.href = "/overview/#{guids[item]}"
    )
    # $('#dc-search .dc-search-input input').keyup () ->
    #   if $(@).val()
    #     $.getJSON "/api/search/#{$(@).val()}", (data) ->
    #       $.get "/templates/dc-search.html/", (template) ->
    #       searchResultsTemplate = Mustache.render $(template), data
    #       $('#ResultsTable').empty().append "<td><tr>Title</tr><tr>Guid</tr></td>"
    #       $('#ResultsTable').append(searchResultsTemplate data[i]) for i in [data.length-1..0]
