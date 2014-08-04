'use strict'

describe 'Service: Charts', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  Charts = {}
  beforeEach inject (_Charts_) ->
    Charts = _Charts_

  it 'should do something', ->
    expect(!!Charts).toBe true
