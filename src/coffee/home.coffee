define ['jquery'], ($) ->
  'use strict'
  $ () ->
    $.getJSON "/api/user", (data) ->
      $('#FillInUserName').text data.Username
      $.getJSON "/api/visited", (data) ->
        if data and data.length
          $('#History').empty().append "<p>Welcome Back, You where last viewing:</p>"
          for entry in data
            do (entry) ->
              entry[1] = "#{entry[1].substring(0,83)}..." if entry[1].length > 85
              $('#History').append "<a href='/charts/#{entry[0]}'>#{entry[1]}</a></br>"
