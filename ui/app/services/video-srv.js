'use strict';

Streaming.service('VideoService', ['$http',
    function($http) {
        this.something = function(query) {
            return $http({
                url: '/something',
                method: "GET",
                params: query
            });
        };
    }
]);