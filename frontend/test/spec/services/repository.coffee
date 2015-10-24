'use strict'

describe 'Service: Repository', ->

  # load the service's module
  beforeEach module 'frontendApp'

  # instantiate service
  Repository = {}
  beforeEach inject (_Repository_) ->
    Repository = _Repository_

  it 'should do something', ->
    expect(!!Repository).toBe true
