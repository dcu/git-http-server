'use strict'

describe 'Service: ServerURL', ->

  # load the service's module
  beforeEach module 'frontendApp'

  # instantiate service
  ServerURL = {}
  beforeEach inject (_ServerURL_) ->
    ServerURL = _ServerURL_

  it 'should do something', ->
    expect(!!ServerURL).toBe true
