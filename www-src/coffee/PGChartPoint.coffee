define ['backbone'], (Backbone) -> 
  class PGChartPoint extends Backbone.Model
	urlRoot: '/api/chartpoints'
	idAttribute: '_id'
	title: null
	text: null
