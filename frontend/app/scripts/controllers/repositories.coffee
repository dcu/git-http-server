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
        Repository.all((response)->
            $scope.repositories = response.data
        ,
            (response)->
                console.log(response)
        )
        return

