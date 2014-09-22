'use strict'

describe 'Service: RequestInterceptor', ->

  # load the service's module
  beforeEach module 'dataplayApp'

  # instantiate service
  RequestInterceptor = {}
  beforeEach inject (_RequestInterceptor_) ->
    RequestInterceptor = _RequestInterceptor_

  it 'should do something', ->
    expect(!!RequestInterceptor).toBe true
