'use strict';

var Streaming = angular.module('Streaming', ['ngRoute', 'ngAnimate', 'ui.bootstrap'])
    .config(function($routeProvider) {
        $routeProvider
            .when('/sample', {
                templateUrl: 'app/views/video.html',
                controller: 'VideoCtrl'
            })
            .otherwise({
                redirectTo: '/sample'
            });
    }).filter('formatDate', function() {
        return function(date) {
            return moment(date).format("YYYY-MM-DD HH:mm:ss.SSS");
        }
    });