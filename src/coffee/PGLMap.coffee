define ['jquery', 'underscore', 'leaflet'], ($, _, L) ->
  class PGLMap
    container: 'body'
    height: '80em'
    map: null
    location: null
    baseLayers: []
    controls: []
    markers: []
    featuresPopupTemplate: _.template $('#features-popup-template').html(), null, variable: 'data'
    externalTrigger: false

    constructor: (container) ->
      @container = container if container
      $(@container).height @height
      @initialize()

    # ---------------------------- Initialization ------------------------------ #
    initialize: ->
      # 0. Create the Map
      @map = L.map @container.substring(1)
      # 1. Add base layers
      @initBaseLayers()
      # 2. Add controls
      @initControls()
      # 3. Register Events
      @registerEvents()
      # Initial map location (geolocation)
      @map.locate setView: true, maxZoom: 16

    # ----------------------------- Base Layers ------------------------------- #
    initBaseLayers: () ->
      baseLayers = [
        {
          name: 'osm'
          layer: (L.tileLayer 'http://{s}.tile.osm.org/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'     
          )
        }
      ]
      @baseLayers.push layer for layer in baseLayers
      layer.layer.addTo @map for layer in @baseLayers

    # ---------------------------- Controls ------------------------------ #
    initControls: ->
      # Layer control
      opts = {}
      opts[layer.name] = layer.layer for layer in @baseLayers

      controls = [
        L.control.layers opts
        #L.control.zoom() it seems to be there by default
      ]
      @controls.push control for control in controls
      control.addTo @map for control in @controls

    # ---------------------------- Events ------------------------------ #
    registerEvents: ->
      @map.on 'locationfound', (e) => @handleGeoLocated e
      @map.on 'locationerror', (e) => @handleGeoLocationError e

      @map.on 'moveend zoomend resize locationfound', (e) => 
        console.log "#{e.type} - #{@externalTrigger}"
        if @externalTrigger
          # Caution!! we're supposing 'moveend' always fires and is also the last
          @externalTrigger = false if e.type is 'moveend' 
        else
          # trigger for updating charts data to map bounds
          $(@).trigger 'update', @map.getBounds()

    # ---------------------------- Geolocation Handlers ------------------------------ #
    handleGeoLocated: (e) -> 
      console.log e
      @location = e.latlng
      radius = e.accuracy / 2
      
      L.circle(e.latlng, radius).addTo(@map)

    handleGeoLocationError: (e) ->
      console.log e 

    updateItems: (items, fitToBounds) ->
      console.log items
      if items and items.length
        @map.removeLayer(marker) for marker in @markers
        @markers = []
        bounds = L.latLngBounds @location, @location
        @addItem item, bounds for item in items
        @map.fitBounds bounds if fitToBounds

    addItem: (item, bounds) ->
      markerLatlng = L.latLng item.lat, item.lon
      marker = L.marker(markerLatlng)
        .addTo(@map)
        .bindPopup(@featuresPopupTemplate item)
        .on("click", (evt) -> marker.openPopup())
      @markers.push marker
      bounds.extend markerLatlng
