<template>
  <div class="container">
    <div class="row">
      <b-table striped hover :items="items"></b-table>
    </div>
    <div class="row">
      <button v-on:click="greet">Greet</button>
    </div>
  </div>
</template>

<script>

export default {
  name: 'instruments',
  data () {
    return {
      items: []
    }
  },
  methods: {
    greet: function (event) {
      this.$socket.send('some data')
    }
  },
  created: function () {
    console.log('created')
    this.$http.get('/instruments/').then(response => {
      this.items = response.data.data

      this.$http.get('/prices/').then(response => {
        this.items.forEach(function (instrument) {
          response.data.data.forEach(function (instrumentPrice) {
            if (instrument.InstrumentKey === instrumentPrice.instrumentKey) {
              instrument.Price = instrumentPrice.price
            }
          })
        })
      })
    })
  },
  mounted () {
    console.log('mounted')
    this.$options.sockets.onmessage = function (message) {
      console.log('Socket message received: ' + message.data)

      var receivedInstrumentPrice = JSON.parse(message.data)

      this.items.forEach(function (instrument) {
        if (instrument.InstrumentKey === receivedInstrumentPrice.instrumentKey) {
          instrument.Price = receivedInstrumentPrice.price
        }
      })
    }
  }
}
</script>

<style>
  .bottom {
    margin-top: 50px;
    margin-left: 200px;
  }

  .bottom-margin {
    margin-bottom: 50px;
  }
</style>
