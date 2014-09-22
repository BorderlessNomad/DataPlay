headlessDeps = ['leafletProviders', 'leafletSearch', 'leafletMarkercluster']
define ['jquery', 'underscore', 'leaflet'].concat(headlessDeps), ($, _, L) ->
  class PGLMap
    container: 'body'
    height: '80em'
    map: null
    location: null
    controls: []
    markers: []
    clusters: []
    featuresPopupTemplate: _.template $('#features-popup-template').html(), null, variable: 'data'
    externalTrigger: false
    baseLayers: []

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
      @map.locate setView: true, maxZoom: 7

    # ----------------------------- Base Layers ------------------------------- #
    initBaseLayers: () ->
      # To define new or redefine existent providers just do it thru leaflet-providers plugin
      # L.TileLayer.Provider.providers.<provider> = ....

    # ---------------------------- Controls ------------------------------ #
    initControls: ->
      # Layers control
      baseLayers = ['OpenStreetMap.Mapnik', 'Stamen.Watercolor']
      overlays = ['OpenWeatherMap.Clouds', 'OpenWeatherMap.Rain']
      layersControl = L.control.layers.provided(baseLayers, overlays)

      # Search control
      searchControl = new L.Control.Search(
        url: 'http://nominatim.openstreetmap.org/search?format=json&q={s}'
        jsonpParam: 'json_callback'
        propertyName: 'display_name'
        propertyLoc: ['lat','lon']
      )

      controls = [
        layersControl
        searchControl
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
        @map.removeLayer(cluster) for cluster in @clusters when cluster
        @clusters = []
        bounds = L.latLngBounds @location, @location
        @addItem item, bounds for item in items
        cluster?.addTo @map for cluster in @clusters
        @map.fitBounds bounds if fitToBounds

    addItem: (item, bounds) ->
      markerLatlng = L.latLng item.lat, item.lon
      marker = L.marker(markerLatlng)       
        .bindPopup(@featuresPopupTemplate item)
        .on("click", (evt) -> marker.openPopup())
      if item.cluster
        @clusters[item.cluster] = new L.MarkerClusterGroup()unless @clusters[item.cluster]
        @clusters[item.cluster].addLayer marker
      else
        marker.addTo @map      
      @markers.push marker
      bounds.extend markerLatlng
