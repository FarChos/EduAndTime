<template>
  <div class="flex flex-row gap-2 pr-28 md:pr-0">

  <img src="../img/inicio/LogoEat.png" alt="" class="h-8 md:h-10">

  <h1 class="text-lg text-white md:text-3xl">{{Text}}</h1>

  </div>
</template>

<script setup lang="ts">
defineProps({
  Text: {
    type: String,
    required: true,
  }
});
</script>