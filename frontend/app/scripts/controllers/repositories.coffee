'use strict'

###*
 # @ngdoc function
 # @name frontendApp.controller:RepositoriesCtrl
 # @description
 # # RepositoriesCtrl
 # Controller of the frontendApp
###
angular.module 'frontendApp'
    .controller 'RepositoriesCtrl', ($scope, Repository)->
        $scope.loadRepositories = ()->
            Repository.all((response)->
                $scope.repositories = response.data
                $scope.selected = $scope.repositories.items[0]
            ,
                (response)->
                    console.log(response)
            )

        $scope.selectRepository = (repository)->
            Repository.find(
                repository.name
            ,
                (response)->
                    $scope.selected = response.data
            ,
                (response)->
                    console.log(response)
            )

        $scope.selected = null;
        $scope.loadRepositories()


        return

