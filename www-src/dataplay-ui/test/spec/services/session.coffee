'use strict'

describe 'Service: Session', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  Session = {}
  beforeEach inject (_Session_) ->
    Session = _Session_

  it 'should do something', ->
    expect(!!Session).toBe true
