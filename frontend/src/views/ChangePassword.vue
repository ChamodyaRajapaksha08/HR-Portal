<template>
  <section class="password-page">
    <header class="page-heading"><div><h2>Change Password</h2><p class="breadcrumb">Dashboard > Change Password</p></div></header>
    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <p v-if="success" class="success-message" role="status">{{ success }}</p>
    <form class="password-card" @submit.prevent="submit">
      <label>Current Password<input v-model="form.current_password" type="password" autocomplete="current-password" /></label>
      <p v-if="errors.current_password" class="field-error">{{ errors.current_password }}</p>
      <label>New Password<input v-model="form.new_password" type="password" autocomplete="new-password" /></label>
      <p v-if="errors.new_password" class="field-error">{{ errors.new_password }}</p>
      <label>Confirm New Password<input v-model="form.confirm_password" type="password" autocomplete="new-password" /></label>
      <p v-if="errors.confirm_password" class="field-error">{{ errors.confirm_password }}</p>
      <button class="primary-button" type="submit" :disabled="saving">{{ saving ? 'Saving...' : 'Update Password' }}</button>
    </form>
  </section>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { apiPost } from '@/api';

const form = reactive({ current_password: '', new_password: '', confirm_password: '' });
const errors = ref({});
const error = ref('');
const success = ref('');
const saving = ref(false);

const submit = async () => {
  errors.value = {};
  error.value = '';
  success.value = '';
  if (!form.current_password) errors.value.current_password = 'Enter your current password.';
  if (form.new_password.length < 8) errors.value.new_password = 'Use at least 8 characters.';
  if (form.new_password !== form.confirm_password) errors.value.confirm_password = 'Passwords do not match.';
  if (Object.keys(errors.value).length) return;

  saving.value = true;
  try {
    const result = await apiPost('/api/profile/change-password', form);
    Object.assign(form, { current_password: '', new_password: '', confirm_password: '' });
    success.value = result.message || 'Password changed successfully.';
  } catch (requestError) {
    error.value = requestError.message || 'Unable to change your password.';
  } finally {
    saving.value = false;
  }
};
</script>

<style scoped>
.password-page { width: 100%; min-width: 0; }
.page-heading { margin-bottom: 20px; }
.page-heading h2 { margin: 0 0 4px; color: #0f172a; font-size: 24px; }
.breadcrumb { margin: 0; color: #64748b; font-size: 13px; }
.password-card { display: grid; gap: 12px; max-width: 520px; padding: 22px; border: 1px solid #e2e8f0; border-radius: 12px; background: white; }
.password-card label { display: grid; gap: 6px; color: #334155; font-size: 13px; font-weight: 600; }
.password-card input { width: 100%; box-sizing: border-box; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 8px; }
.primary-button { justify-self: start; min-height: 38px; padding: 8px 14px; border: 0; border-radius: 8px; background: #047857; color: white; font-size: 13px; font-weight: 600; cursor: pointer; }
.field-error { margin: -8px 0 0; color: #b91c1c; font-size: 12px; }
.api-error { color: #b91c1c; }
.success-message { color: #047857; }
</style>
