define ['jquery', 'underscore', 'backbone'], ($, _, Backbone) -> 
  class PGChartPointView extends Backbone.View
    tagName: 'div'
    className: 'pointDialog'
    template: _.template $('#pointDataTemplate').html(), null, variable: 'data'

    events:
      "click .pointSave" : "savePoint"
      "click .pointRemove" : "removePoint"
      "click .pointCancel" : "closeDialog"

    initialize: ->
    	@listenTo(@model, 'change', @render)

    render: ->
    	@$el.html(@template(@model.attributes))
    	@

    savePoint: ->
    	alert "Saving point #{@model.get('title')}"

    removePoint: ->
    	alert "Removing point #{@model.get('title')}"
    	@remove()

    closeDialog: -> @remove()