'use strict'

###*
 # @ngdoc function
 # @name frontendApp.controller:DetailsCtrl
 # @description
 # # DetailsCtrl
 # Controller of the frontendApp
###
angular.module 'frontendApp'
    .controller 'DetailsCtrl', ($scope, $routeParams, $showdown, Repository)->
        Repository.find(
            $routeParams.path
        ,
            (response) ->
                $scope.details = response.data
                $scope.details.readme_html = $showdown.makeHtml($scope.details.readme)
        ,
            (response) ->
                console.log(response)
        )
        return
