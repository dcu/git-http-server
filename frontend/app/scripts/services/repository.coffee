'use strict'

###*
 # @ngdoc service
 # @name frontendApp.Repository
 # @description
 # # Repository
 # Factory in the frontendApp.
###
angular.module 'frontendApp'
  .factory 'Repository', ($http, $showdown, ApiConfig) ->

    # Public API here
    all: (successCallback, errorCallback) ->
        $http({method: "GET", url: "#{ApiConfig.url}/repositories"}).then(successCallback, errorCallback)

    find: (path, successCallback, errorCallback) ->
        $http({method: "GET", url: "#{ApiConfig.url}/repositories/#{path}"}).then(
            (response) -> 
                response.data.readme_html = $showdown.makeHtml(response.data.readme)
                successCallback(response)
        , 
            errorCallback)


