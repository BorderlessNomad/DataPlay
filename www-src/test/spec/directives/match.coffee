'use strict'

describe 'Directive: match', ->

  # load the directive's module
  beforeEach module 'dataplayApp'

  scope = {}

  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()

  it 'should make hidden element visible', inject ($compile) ->
    element = angular.element '<match></match>'
    element = $compile(element) scope
    expect(element.text()).toBe 'this is the match directive'
