<template>
  <div class="container">
    <div class="row">
      <b-table striped hover :items="items"></b-table>
    </div>
  </div>
</template>

<script>
export default {
  name: 'searchorder',
  data () {
    return {
      query: '',
      items: []
    }
  },
  methods: {
    search () {
      this.$http.get('/orders/search?q=' + this.query)
        .then(response => {
          this.items = response.data.data
        })
    }
  },
  created: function () {
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
