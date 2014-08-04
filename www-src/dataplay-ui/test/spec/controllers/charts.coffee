'use strict'

describe 'Controller: ChartsCtrl', ->

  # load the controller's module
  beforeEach module 'dataplayApp'

  ChartsCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    ChartsCtrl = $controller 'ChartsCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
