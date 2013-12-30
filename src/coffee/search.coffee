define ['jquery'], ($) ->
  'use strict'
  $ () ->
    guid = window.location.href.split('/')[window.location.href.split('/').length - 1]
    $('#SBox').keyup () ->
      if $('#SBox').val()
        $.getJSON "/api/search/#{$('#SBox').val()}", (data) ->
          $('#ResultsTable').empty().append "<td><tr>Title</tr><tr>Guid</tr></td>"
          for i in [data.length-1..0]
            do (i) ->
              if data[i].Title
                $('#ResultsTable').append(
                  "<tr id='#{data[i].GUID}'>" + 
                    "<td style='width: 70%;'>#{data[i].Title}</td>" + 
                    "<td id='Note#{data[i].GUID }'>#{data[i].GUID}</td>" + 
                  "</tr>"
                )
                $("##{data[i].GUID}").addClass("success")
                if guid isnt "overlay"
                  $("#Note#{data[i].GUID}").html(
                    "<a href='/charts/#{data[i].GUID}'>View &raquo;</a>&nbsp;" +
                    "<a href='/grid/#{data[i].GUID}'>Grid &raquo;</a>&nbsp;" + 
                    "<a href='/overview/#{data[i].GUID}'>Overview &raquo;</a>&nbsp;"
                  )
                else
                  $('#Note' + data[i].GUID).html "<a href='/overlay/#{data[i].GUID}'>Overlay &raquo;</a>"
                  console.log data

