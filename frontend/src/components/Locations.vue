<template>
  <div id="locations_root">
    <p id="update_info">Last location update: {{ lastUpdate }}</p>

    <div id="map"></div>
    
  </div>
</template>

<script>
import gmapsInit from '../plugins/gmaps';

export default {
  async mounted() {
    try {
      const google = await gmapsInit();
      const geocoder = new google.maps.Geocoder();
      const map = new google.maps.Map(document.getElementById("map"));

      geocoder.geocode({ address: 'Germany' }, (results, status) => {
        if (status !== 'OK' || !results[0]) {
          throw new Error(status);
        }

        map.setCenter(results[0].geometry.location);
        map.fitBounds(results[0].geometry.viewport);
      });
    } catch (error) {
      console.error(error);
    }
  },
  data () {
    return {
      lastUpdate: "2019-04-03 11:55:02",
      map: null,
    }
  }
} 

</script>

<style>
#map {
height: 100%;
}
#update_info {
  padding: 20px;
  margin: 0px;
}
.container, #locations_root {
  height: 95%;
  margin: 0;
  padding: 0;
}
</style>