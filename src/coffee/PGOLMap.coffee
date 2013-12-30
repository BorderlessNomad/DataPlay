define ['jquery', 'openlayers'], ($, OpenLayers) ->
  OpenLayers.ImgPath = "/img/openlayers/"
  OpenLayers.ProxyHost = "proxy.cgi?url="

  class PGOLMap
    @OSM_PROJECTION: new OpenLayers.Projection("EPSG:4326")
    @OL_PROJECTION: new OpenLayers.Projection("EPSG:900913")
    container: 'body'
    width: '60em'
    height: '30em'
    readonly: false
    map: null
    controls: []
    baseLayers: []
    geoLocationControl: null
    searchBarTemplate: _.template $('#search-bar-template').html()
    locationFeature: null
    searchBounds: null
    searchScope: 0.05

    constructor: (container) ->
      @container = container if container
      @width = $(@container).width()
      @height = $(@container).height()
      @initialize()

    # ---------------------------- Initialization ------------------------------ #
    initialize: ->
      # 0. Create the Map
      @map = new OpenLayers.Map(
        @container
        theme: null
        projection: "EPSG:900913"
        fractionalZoom: true
      )
      # 1. Add base layers
      @initBaseLayers()
      # 2. Add controls
      @initControls()
      # Initial map zoom
      @map.zoomToMaxExtent()

    # ---------------------------- Renderization ------------------------------ #
    render: -> @

    # ----------------------------- Base Layers ------------------------------- #
    initBaseLayers: () ->
      baseLayers = [
        new OpenLayers.Layer.OSM(
          "OpenStreetMap"
          [ "http://a.tile.openstreetmap.org/${z}/${x}/${y}.png"
            "http://b.tile.openstreetmap.org/${z}/${x}/${y}.png"
            "http://c.tile.openstreetmap.org/${z}/${x}/${y}.png" ]
          isBaseLayer: true
          transitionEffect: 'resize'
        )
        new OpenLayers.Layer.Bing(
          name: "Bing Road"
          key: "AqTGBsziZHIJYYxgivLBf0hVdrAk9mWO5cQcb8Yux8sW5M8c8opEC2lZqKR1ZZXf"
          type: "Road"
          metadataParams: { mapVersion: "v1" }
          isBaseLayer: true
        )
        new OpenLayers.Layer.Bing(
          name: "Bing Aerial"
          key: "AqTGBsziZHIJYYxgivLBf0hVdrAk9mWO5cQcb8Yux8sW5M8c8opEC2lZqKR1ZZXf"
          type: "Aerial"
          isBaseLayer: true
        )
        new OpenLayers.Layer.Bing(
          name: "Bing Aerial With Labels"
          key: "AqTGBsziZHIJYYxgivLBf0hVdrAk9mWO5cQcb8Yux8sW5M8c8opEC2lZqKR1ZZXf"
          type: "AerialWithLabels"
          isBaseLayer: true
        )
      ]
      @baseLayers.push lyr for lyr in baseLayers
      @map.addLayers @baseLayers

      # Markers layer
      @featuresLayer = new OpenLayers.Layer.Markers "Features" # change this to features layer
      @map.addLayer @featuresLayer

    # ---------------------------- Controls ------------------------------ #
    initControls: ->
      # Cache controls
      cacheRead = new OpenLayers.Control.CacheRead()
      @controls.push cacheRead

      cacheWrite = new OpenLayers.Control.CacheWrite
        autoActivate: true
        imageFormat: "image/jpeg"
        eventListeners:
          cachefull: () -> console.log "Cache Full"
      @controls.push cacheWrite
      # Reset the cache on init
      OpenLayers.Control.CacheWrite.clearCache()

      # Custom Keyboard control & handler
      keyboardControl = new OpenLayers.Control
      keyboardControl.handler = new OpenLayers.Handler.Keyboard(
        keyboardControl
        'keyup': (e) => dummy = true #console.log "key #{e.keyCode}"
        # TODO: handle different key events
      )
      @controls.push keyboardControl

      # Geolocation layer & control
      geoLocationLayer = new OpenLayers.Layer.Vector("Your location")
      @map.addLayer(geoLocationLayer);
      @geoLocationControl = new OpenLayers.Control.Geolocate
        bind: true
        watch: true
        geolocationOptions:
          enableHighAccuracy: true
          maximumAge: 0
          timeout: 7000
      @geoLocationControl.follow = true
      @geoLocationControl.events.register(
        "locationupdated"
        @
        (e) ->
          geoLocationLayer.removeAllFeatures()
          @locationFeature = new OpenLayers.Feature.Vector(
            e.point
            {}
            {
              graphicName: 'circle'
              strokeColor: '#0000FF'
              strokeWidth: 1
              fillOpacity: 0.5
              fillColor: '#0000BB'
              pointRadius: 20
            }
          )
          geoLocationLayer.addFeatures [@locationFeature]
          @handleGeoLocated e
      )
      @geoLocationControl.events.register(
        "locationfailed"
        @
        () -> OpenLayers.Console.log 'Location detection failed'
      )
      @controls.push @geoLocationControl

      # Search control
      that = @
      OpenLayers.Control.prototype.keepEvents = (div) ->
        @keepEventsDiv = new OpenLayers.Events(@, div, null, true)

        triggerSearch = (evt) =>
          element = OpenLayers.Event.element(evt)
          if  evt.keyCode == 13 then that.performSearch $(element).val()

        @keepEventsDiv.on
          "mousedown": (evt) ->
            @mousedown = true
            OpenLayers.Event.stop(evt, true)
          "mousemove": (evt) ->
            OpenLayers.Event.stop(evt, true) if @mousedown
          "mouseup": (evt) ->
            if @mousedown
              @mousedown = false
              OpenLayers.Event.stop(evt, true)
          "click": (evt) -> OpenLayers.Event.stop(evt, true)
          "mouseout": (evt) -> @mousedown = false
          "dblclick": (evt) -> OpenLayers.Event.stop(evt, true)
          "touchstart": (evt) -> OpenLayers.Event.stop(evt, true)
          "keydown": (evt) -> triggerSearch(evt)
          scope: @

      searchControl = new OpenLayers.Control
      OpenLayers.Util.extend searchControl,
        displayClass: 'searchControl'
        initialize : () ->
          OpenLayers.Control.prototype.initialize.apply(@, arguments)
        draw: () ->
          div = OpenLayers.Control.prototype.draw.apply(@, arguments)
          div.innerHTML = that.searchBarTemplate {}
          @keepEvents(div);
          $(div).find('button.geosearchButton').click () =>
            that.performSearch $('#searchInput').val()
          $(div).find('button.geoLocateButton').click () =>
            that.geoLocationControl.deactivate()
            that.geoLocationControl.activate()
          div
        allowSelection: true
      @controls.push searchControl

      # Simple map position display control
      mapPosition = new OpenLayers.Control.MousePosition()
      mapPosition.displayProjection = new OpenLayers.Projection("EPSG:4326")
      @controls.push mapPosition

      # Layer switcher control
      @controls.push new OpenLayers.Control.LayerSwitcher()

      # TODO: Add more controls here ....

      @map.addControls @controls

      # Controls activation
      keyboardControl.activate()
      @geoLocationControl.activate()
      searchControl.activate()

    # ---------------------------- Geolocation Handler ------------------------------ #
    handleGeoLocated: (e) -> 
      console.log e
      r = @searchScope
      p = new OpenLayers.Geometry.Point e.point.x, e.point.y
      p.transform @map.getProjectionObject(), Maps.MapView.OSM_PROJECTION
      @searchBounds = [p.y - r, p.x - r, p.y + r, p.x + r]
      @map.zoomToExtent e.point.getBounds()

    # ---------------------------- Search handler ------------------------------ #
    performSearch: (text) ->
      $.get "http://nominatim.openstreetmap.org/search?q=#{text}&format=json&limit=10", (data) =>
        if data? and data.length > 0
          $('div.searchControl ul.searchList').css('display', 'block').html('')
          for item in data
            do (item) =>
              $('div.searchControl ul.searchList').append "<li>#{item.display_name.substring(0, 20)}</li>"
              $('div.searchControl ul.searchList').find('li').last().click () =>
                $('div.searchControl ul.searchList').css('display', 'none').html('')
                @selectLocation item

    selectLocation: (location) -> 
      console.log location
      p = new OpenLayers.Geometry.Point location.lon, location.lat
      p.transform Maps.MapView.OSM_PROJECTION, @map.getProjectionObject()
      @locationFeature = new OpenLayers.Feature.Vector(
        p
        {}
        externalGraphic: location.icon
        pointRadius: 10
      )
      @drawLayer.removeAllFeatures()
      @drawLayer.addFeatures [@locationFeature]
      @map.zoomToExtent p.getBounds()



