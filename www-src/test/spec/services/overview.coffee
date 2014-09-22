'use strict'

describe 'Service: Overview', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  Overview = {}
  beforeEach inject (_Overview_) ->
    Overview = _Overview_

  it 'should do something', ->
    expect(!!Overview).toBe true
