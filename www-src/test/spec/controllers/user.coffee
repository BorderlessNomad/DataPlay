'use strict'

describe 'Controller: UserCtrl', ->

  # load the controller's module
  beforeEach module 'dataplayApp'

  UserCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    UserCtrl = $controller 'UserCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
