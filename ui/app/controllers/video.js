'use strict';

angular.module('Streaming').controller('VideoCtrl', function($scope, $location, $routeParams, VideoService) {

    $scope.something = function() {

        VideoService.something($scope.query)
            .success(function(response) {

            })
            .error(function() {

            });
    };
});
