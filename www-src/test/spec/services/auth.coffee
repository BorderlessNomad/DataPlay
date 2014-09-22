'use strict'

describe 'Service: Auth', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  Auth = {}
  beforeEach inject (_Auth_) ->
    Auth = _Auth_

  it 'should do something', ->
    expect(!!Auth).toBe true
