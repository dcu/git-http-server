'use strict'

###*
 # @ngdoc function
 # @name frontendApp.controller:DetailsCtrl
 # @description
 # # DetailsCtrl
 # Controller of the frontendApp
###
angular.module 'frontendApp'
  .controller 'DetailsCtrl', ($scope, $routeParams, Repository)->
    Repository.find(
        $routeParams.path
    ,
        (response) ->
            $scope.details = response.data
    ,
        (response) ->
            console.log(response)
    )
    return
