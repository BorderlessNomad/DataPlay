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
          name: 'OpenStreetMaps(default)'
          layer: (L.tileLayer 'http://{s}.tile.osm.org/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'     
          )
        }
        {
          name: 'OpenStreetMaps(B&W)'
          layer: (L.tileLayer 'http://{s}.www.toolserver.org/tiles/bw-mapnik/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'     
          )
        }
        {
          name: 'OpenStreetMaps(Deutch)'
          layer: (L.tileLayer 'http://{s}.tile.openstreetmap.de/tiles/osmde/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'     
          )
        }
        {
          name: 'OpenStreetMaps(Hot)'
          layer: (L.tileLayer 'http://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png',
            attribution: '&copy; Tiles courtesy of <a href="http://hot.openstreetmap.org/" target="_blank">Humanitarian OpenStreetMap Team</a>'     
          )
        }
        {
          name: 'OpenCycleMap'
          layer: (L.tileLayer 'http://{s}.tile.opencyclemap.org/cycle/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://www.opencyclemap.org">OpenCycleMap</a>'     
          )
        }
        # {
        #   name: 'OpenSeaMap'
        #   layer: (L.tileLayer 'http://tiles.openseamap.org/seamark/{z}/{x}/{y}.png',
        #     attribution: 'Map data: &copy; <a href="http://www.openseamap.org">OpenSeaMap</a> contributors'     
        #   )
        # }

        {
          name: 'Bing(Aerial)'
          layer: (L.tileLayer 'http://{s}.tile.opencyclemap.org/cycle/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://www.opencyclemap.org">OpenCycleMap</a>'     
          )
        }

        {
          name: 'Thunderforest(Transport)'
          layer: (L.tileLayer 'http://{s}.tile2.opencyclemap.org/transport/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://www.opencyclemap.org">OpenCycleMap</a>'     
          )
        }
        {
          name: 'Thunderforest(Landscape)'
          layer: (L.tileLayer 'http://{s}.tile3.opencyclemap.org/landscape/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://www.opencyclemap.org">OpenCycleMap</a>'     
          )
        }
        {
          name: 'Thunderforest(Outdoors)'
          layer: (L.tileLayer 'http://{s}.tile.thunderforest.com/outdoors/{z}/{x}/{y}.png',
            attribution: '&copy; <a href="http://www.opencyclemap.org">OpenCycleMap</a>'     
          )
        }

        {
          name: 'Stamen(Default)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(TonerBackground)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner-background/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(TonerHybrid)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner-hybrid/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(TonerLines)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner-lines/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(TonerLabels)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner-labels/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(Lite)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/toner-lite/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(Terrain)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/terrain/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(TerrainBackground)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/terrain-background/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
          )
        }
        {
          name: 'Stamen(Watercolor)'
          layer: (L.tileLayer 'http://{s}.tile.stamen.com/watercolor/{z}/{x}/{y}.png',
            attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash;'     
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
