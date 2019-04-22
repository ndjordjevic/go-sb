<template>
  <div class="container py-3">
    <div class="row">
      <div class="mx-auto col-sm-6">
        <div class="card">
          <div class="card-header">
            <h4 class="mb-0">Enter new order</h4>
          </div>
          <div class="card-body">
            <form @submit.prevent="createNewOrder" autocomplete="off" class="form" role="form">
              <div class="form-group row">
                <label class="col-lg-4 col-form-label form-control-label">Email</label>
                <div class="col-lg-8">
                  <input class="form-control" placeholder="Enter email" type="text" v-model="order.email">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-lg-4 col-form-label form-control-label">Instrument Key</label>
                <div class="col-lg-8">
                  <input class="form-control" placeholder="Instrument Key" type="text" v-model="order.instrumentKey">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-lg-4 col-form-label form-control-label">Currency</label>
                <div class="col-lg-8">
                  <input class="form-control" placeholder="Currency" type="text" v-model="order.currency">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-lg-4 col-form-label form-control-label">Size</label>
                <div class="col-lg-8">
                  <input class="form-control" placeholder="Size" type="number" v-model="order.size">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-lg-4 col-form-label form-control-label">Price</label>
                <div class="col-lg-8">
                  <input class="form-control" placeholder="Price" type="number" v-model="order.price">
                </div>
              </div>
              <div class="form-group row">
                <label class="col-lg-3 col-form-label form-control-label"></label>
                <div class="col-lg-9">
                  <input class="btn btn-primary" type="submit" value="Save Order">
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  // @ is an alias to /src
  // import HelloWorld from '@/components/HelloWorld.vue'

  export default {
    name: 'order',
    components: {
      // HelloWorld
    },
    data() {
      return {
        order: {
          email: '',
          instrumentKey: '',
          currency: '',
          size: '',
          price: ''
        }
      }
    },
    methods: {
      createNewOrder() {
        console.log(this.order)
      }
    }
  }
</script>
