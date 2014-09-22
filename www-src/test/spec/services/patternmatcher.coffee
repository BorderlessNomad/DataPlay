'use strict'

describe 'Service: PatternMatcher', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  PatternMatcher = {}
  beforeEach inject (_PatternMatcher_) ->
    PatternMatcher = _PatternMatcher_

  it 'should do something', ->
    expect(!!PatternMatcher).toBe true
