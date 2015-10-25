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
                $scope.details.repo_url = "#{window.location.protocol}//#{window.location.host}/#{$scope.details.repository.name}.git"
        ,
            (response) ->
                console.log(response)
        )
        return
