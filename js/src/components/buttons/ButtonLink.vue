<template>
  <div>

    <router-link v-if="to" :class="cls" :to="to">
      <div :class="'icon ' + iconClass">
        <icon :name="icon"></icon>
      </div>
      <p>{{label}}</p>
      <div class="clear"></div>
    </router-link>

    <a v-else :class="cls" @click="func">
      <div :class="'icon ' + iconClass">
        <icon :name="icon"></icon>
      </div>
      <p>{{label}}</p>
      <div class="clear"></div>
    </a>

  </div>
</template>

<script>
export default {
  name: "ButtonLink",
  // TODO(thiderman): Implement tooltip
  props: ['icon', 'iconClass', 'label', 'to', 'cls', 'func', 'tooltip'],
}
</script>
