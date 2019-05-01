<template>
  <div class="container">
    <div class="row bottom-margin">
      <div class="col">
        <div class="input-group input-group-lg bottom">
          <div class="input-group-prepend">
            <span class="input-group-text">Search</span>
          </div>
          <input type="text"
                 class="form-control col-md-6"
                 @keyup.prevent="search"
                 v-model="query"/>
        </div>
      </div>
    </div>
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
