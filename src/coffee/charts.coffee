define [
  'jquery'
  'underscore'
  'd3'
  'app/Common'
  'app/PGPatternMatcher'
  'app/PGLinesChart'
  'app/PGAreasChart'
  'app/PGBarsChart'
  'app/PGPieChart'
  'app/PGBubblesChart'
  'app/PGEnclosureChart'
  'app/PGTreeChart'
  'app/PGTreemapChart'
], (
  $
  _
  d3
  Common
  PGPatternMatcher
  PGLinesChart
  PGAreasChart
  PGBarsChart
  PGPieChart
  PGBubblesChart
  PGEnclosureChart
  PGTreeChart
  PGTreemapChart
) ->
  'use strict'
  DataCon = {}
  aUrl = window.location.href.split('/')
  last = aUrl[aUrl.length-1]
  aGuid = last.split('?')
  guid = aGuid[0]
  qs = if aGuid.length>1 then decodeURI(aGuid[1]) else null

  go2HierarchyChart = ->
    count = _.keys(DataCon.dataset[0]).length
    selectorTemplate = _.template($('#selectorTemplate').html(), null, variable: 'data')
    selectorButton = $('<a class="btn btn-lg btn-success" href="javascript:void(0)" role="button">Draw Chart</a>')
    data = type: 'Key', keys: []
    dataset = key: '', values: null             
    data.keys.push(node) for node of DataCon.dataset[0]
    $('#selectors').html ''
    $('#selectors').append(selectorTemplate data) for i in [0..count-2] 
    data.type = 'Value'
    $('#selectors').append(selectorTemplate data)
    $('#controls').html('').append(selectorButton)

    selectorButton.click ->
      $("#chart").html ''
      keySelectors = $('#selectors select.Key')
      valueSelector = $('#selectors select.Value')
      nested = d3.nest()
      for i in [0..keySelectors.length-1]
        do (i) ->
          val = $(keySelectors.get(i)).val()
          if val and val isnt '0'
            ((theNest, theNode) -> theNest.key (d) -> d[theNode])(nested, $(keySelectors.get(i)).val())
      dataCopy = JSON.parse(JSON.stringify DataCon.dataset)
      nested = nested.entries dataCopy
      dataset.values = nested
      switch DataCon.currChartType
        when 'enclosure'
          DataCon.chart = new PGEnclosureChart("#chart", {top: 0, right: 0, bottom: 0, left: 0}, dataset, null, DataCon.patterns, null, valueSelector.val())
        when 'tree'
          DataCon.chart = new PGTreeChart("#chart", {top: 0, right: 0, bottom: 0, left: 0}, dataset, null, DataCon.patterns, null, valueSelector.val())
        when 'treemap'
          DataCon.chart = new PGTreemapChart("#chart", {top: 0, right: 0, bottom: 0, left: 0}, dataset, null, DataCon.patterns, null, valueSelector.val())
    
    if DataCon.autokeys and DataCon.autokeys.length
      selects = $('#selectors select')
      $(selects.get(i)).val(DataCon.autokeys[i]) for i in [0..DataCon.autokeys.length-2]
      $(selects[selects.length-1]).val(DataCon.autokeys[DataCon.autokeys.length-1])
      selectorButton.click()

  go2Chart = (type) ->
    if DataCon.currChartType isnt type
      DataCon.currChartType = type
      DataCon.chart = null
      $('#modes').html ''
      defaultSelectorTemplate = _.template $('#defaultSelectorTemplate').html()
      $('#controls').html ''
      $('#selectors').html(defaultSelectorTemplate {})
      $('#SetupOverlay').click ->
        #Okay so we need to put into local storage the current GUID so when 
        #the overlay panal is called we can use it later on the overlay page.
        localStorage['overlay1'] = guid
        window.location.href = '/search/overlay'
      populateKeys()

  updateChart = ->
    chartAxes = x: $("#pickxaxis").val(), y: $("#pickyaxis").val()
    chartData = Common.parseChartData(
      DataCon.dataset
      {key: chartAxes.x, pattern: DataCon.patterns[chartAxes.x]}
      {key: chartAxes.y, pattern: DataCon.patterns[chartAxes.y]}
    )   
    console.log chartData
    if not DataCon.chart
      $("#chart").html('');
      if qs
        params = qs.split('&')
        entries = {}
        size = 0;
        for param in params
          do (param) ->
            aParam = param.split '='
            entries[aParam[0]] = aParam[1]
            size++
        if entries['chart']
          DataCon.currChartType = entries['chart']
          if ['enclosure', 'tree', 'treemap'].indexOf(entries['chart'])>=0
            DataCon.autokeys = []
            DataCon.autokeys.push(entries["key#{i}"]) for i in [0..size-3]
            DataCon.autokeys.push(entries['value']);
          else 
            if entries['x']
              $("#pickxaxis").val entries['x']
              chartAxes.x = entries['x']
            if entries['y']
              $("#pickyaxis").val entries['y']
              chartAxes.y = entries['y']
            chartData = Common.parseChartData(
              DataCon.dataset
              {key: chartAxes.x, pattern: DataCon.patterns[chartAxes.x]}
              {key: chartAxes.y, pattern: DataCon.patterns[chartAxes.y]}
            )
        qs = null

      switch DataCon.currChartType
        when 'bars'
          DataCon.chart = new PGBarsChart("#chart", null, chartData, chartAxes, DataCon.patterns, null)
        when 'area'
          DataCon.chart = new PGAreasChart("#chart", null, chartData, chartAxes, DataCon.patterns, null)
        when 'pie'
          DataCon.chart = new PGPieChart("#chart", null, chartData, chartAxes, DataCon.patterns, null)
        when 'bubbles'
          dataset = [];
          dataset.push {name: d[0], count: d[1]} for d in chartData
          DataCon.chart = new PGBubblesChart("#chart", {top: 5, right: 0, bottom: 0, left: 0}, dataset, null, DataCon.patterns, null);
        when 'enclosure', 'tree', 'treemap'
          go2HierarchyChart();
        else
          DataCon.chart = new PGLinesChart("#chart", null, chartData, chartAxes, DataCon.patterns, null);        
    else
      switch DataCon.currChartType
        when 'enclosure'
          #Special treatment????
        else 
          DataCon.chart.updateChart(chartData, chartAxes);
    $('#pickxaxis, #pickyaxis').change handleAxesChange

  populateKeys = ->
    $("#pickxaxis").empty()
    $("#pickyaxis").empty()
    $("#pickxaxis").append $("<option></option>").val(key).text(key) for key in DataCon.keys
    $("#pickyaxis").append $("<option></option>").val(key).text(key) for key in DataCon.keys
    $("#pickyaxis").val(DataCon.keys[1]) if DataCon.keys.length > 1
    if qs
      updateChart()
    else
      Common.getUserDefaults(guid, $("#pickxaxis"), $("#pickyaxis"), updateChart)    

  handleAxesChange = ->
    Common.saveUserDefaults guid, $("#pickxaxis").val(), $("#pickyaxis").val()
    updateChart()


  LightUpBookmarks = -> #Do nothing

  Annotations = [];

  SavePoint = (x, y) ->
    SaveObj =
        xax: $("#pickxaxis").val()
        yax: $("#pickyaxis").val()
        x: x
        y: y
        guid: guid
    Annotations.push SaveObj
    #Save it as well and rewrite the titlebar URL to be the shareable one.
    $.ajax(
        type: "POST"
        url: '/api/setbookmark/'
        data: "data=" + JSON.stringify(Annotations)
        success: (resp) -> window.history.pushState 'page2', 'Title', "/viewbookmark/#{resp}" #Saved and sound
        error: (err) -> console.log(err);
    )

  $ ->
    $('.chartSelector').click () -> go2Chart $(@).attr('chart')

    $.getJSON "/api/getinfo/#{guid}", (data) ->
      $('#FillInDataSet').html data.Title
      $("#wikidata").html data.Notes
    
    $.getJSON "/api/getreduceddata/#{guid}/10/100", (data) ->
      console.log(data);
      #so data is an array of shit.
      DataCon.dataset = data
      DataCon.keys = []
      if data and data.length
        DataCon.patterns = {}
        for key of data[0]
          do (key) ->
            DataCon.keys.push key
            # Frontend Pattern Recognition
            DataCon.patterns[key] = {
              valuePattern: PGPatternMatcher.getPattern data[0][key]
              keyPattern: PGPatternMatcher.getKeyPattern data[0][key]
            }

        $.ajax(
          async: false
          url: "/api/identifydata/#{guid}"
          dataType: "json"
          success: (data) ->
            Cols = data.Cols
            for col in Cols
              do (col) ->
                switch col.Sqltype
                  when "int", "bigint"
                    DataCon.patterns[col.Name].valuePattern = 'intNumber'
                  when "float"
                    DataCon.patterns[col.Name].valuePattern = 'floatNumber'
                  else
                    #leave pattern as it was recognised by frontend
            go2Chart 'enclosure'
        )

    $(window).resize ->
      DataCon.chart = null
      updateChart()
