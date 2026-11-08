<template>
  <div class="empty-state" :style="{ minHeight: height }" role="status">
    {{ message }}
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: String,
    default: 'No data yet',
  },
  height: {
    type: String,
    default: 'auto',
  },
});
</script>

<style scoped>
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 18px;
  color: #64748b;
  text-align: center;
  font-size: 13px;
}
</style>
