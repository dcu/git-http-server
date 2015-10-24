'use strict'

###*
 # @ngdoc function
 # @name frontendApp.controller:AboutCtrl
 # @description
 # # AboutCtrl
 # Controller of the frontendApp
###
angular.module 'frontendApp'
  .controller 'AboutCtrl', ->
    @awesomeThings = [
      'HTML5 Boilerplate'
      'AngularJS'
      'Karma'
    ]
    return
