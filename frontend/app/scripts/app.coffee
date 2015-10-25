'use strict'

###*
 # @ngdoc overview
 # @name frontendApp
 # @description
 # # frontendApp
 #
 # Main module of the application.
###
angular
  .module 'frontendApp', [
    'ngAnimate',
    'ngCookies',
    'ngMessages',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch',
    'ng-showdown'
  ]
  .config ($routeProvider) ->
    $routeProvider
      .when '/',
        templateUrl: 'views/repositories.html'
        controller: 'RepositoriesCtrl'
        controllerAs: 'repositories'
      .when '/r/:path*',
        templateUrl: 'views/details.html'
        controller: 'DetailsCtrl'
        controllerAs: 'details'
      .when '/about',
        templateUrl: 'views/about.html'
        controller: 'AboutCtrl'
        controllerAs: 'about'
      .otherwise
        redirectTo: '/'

