'use strict'

describe 'Controller: DetailsCtrl', ->

  # load the controller's module
  beforeEach module 'frontendApp'

  DetailsCtrl = {}

  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    DetailsCtrl = $controller 'DetailsCtrl', {
      # place here mocked dependencies
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(DetailsCtrl.awesomeThings.length).toBe 3
