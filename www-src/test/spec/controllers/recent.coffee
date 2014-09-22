'use strict'

describe 'Controller: RecentCtrl', ->

  # load the controller's module
  beforeEach module 'dataplayApp'

  RecentCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    RecentCtrl = $controller 'RecentCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
