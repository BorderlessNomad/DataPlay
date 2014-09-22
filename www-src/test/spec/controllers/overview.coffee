'use strict'

describe 'Controller: OverviewCtrl', ->

  # load the controller's module
  beforeEach module 'dataplayApp'

  OverviewCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    OverviewCtrl = $controller 'OverviewCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
