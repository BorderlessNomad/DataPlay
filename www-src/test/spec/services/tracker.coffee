'use strict'

describe 'Service: Tracker', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  Tracker = {}
  beforeEach inject (_Tracker_) ->
    Tracker = _Tracker_

  it 'should do something', ->
    expect(!!Tracker).toBe true
