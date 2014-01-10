define ['jquery', 'mustache'], ($, Mustache) ->
  'use strict'
  $ () ->    
    input = $('#dc-search .dc-search-input input')
    res = $('#dc-search .dc-search-input .dc-search-results')
    input.keyup () ->
      if $(@).val()
        $.getJSON "/api/search/#{$(@).val()}", (data) ->
          $.get "/templates/dc-search-results.html/", (template) ->
            fixedData = (item for item in data when item.Title)
            res.fadeIn(200) if fixedData
            res.empty().append Mustache.render(template, data: fixedData)
      else
        res.fadeOut(200)
