'use strict'

###*
 # @ngdoc service
 # @name frontendApp.Repository
 # @description
 # # Repository
 # Factory in the frontendApp.
###
angular.module 'frontendApp'
  .factory 'Repository', ($http, ServerURL) ->

    # Public API here
    all: (successCallback, errorCallback) ->
        $http({method: "GET", url: "#{ServerURL}/repositories"}).then(successCallback, errorCallback)

    find: (path, successCallback, errorCallback) ->
        $http({method: "GET", url: "#{ServerURL}/repositories/#{path}"}).then(successCallback, errorCallback)


