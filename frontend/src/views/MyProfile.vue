<template>
  <section class="self-service-page">
    <header class="page-heading">
      <div><h2>My Profile</h2><p class="breadcrumb">Dashboard > My Profile</p></div>
      <router-link class="secondary-button" to="/dashboard/change-password">Change Password</router-link>
    </header>

    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <p v-if="success" class="success-message" role="status">{{ success }}</p>
    <EmptyState v-if="loading" message="Loading profile..." />
    <EmptyState v-else-if="!profile" message="Profile data is unavailable" />
    <section v-else class="profile-card">
      <template v-if="!editing">
        <dl class="profile-details">
          <div><dt>Name</dt><dd>{{ profile.name }}</dd></div>
          <div><dt>Employee ID</dt><dd>{{ profile.employee_id }}</dd></div>
          <div><dt>Email</dt><dd>{{ profile.email }}</dd></div>
          <div><dt>Position</dt><dd>{{ profile.position || '—' }}</dd></div>
          <div><dt>Department</dt><dd>{{ profile.department || '—' }}</dd></div>
          <div><dt>Role</dt><dd>{{ profile.role }}</dd></div>
        </dl>
        <button class="primary-button" type="button" @click="beginEdit">Edit Profile</button>
      </template>
      <form v-else class="profile-form" @submit.prevent="saveProfile">
        <label>Name<input v-model.trim="form.name" type="text" required /></label>
        <p v-if="fieldError" class="field-error">{{ fieldError }}</p>
        <label>Email<input :value="profile.email" type="email" disabled /></label>
        <label>Position<input v-model.trim="form.position" type="text" /></label>
        <div class="form-actions">
          <button class="primary-button" type="submit" :disabled="saving">{{ saving ? 'Saving...' : 'Save Profile' }}</button>
          <button class="secondary-button" type="button" @click="editing = false">Cancel</button>
        </div>
      </form>
    </section>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { apiGet, apiPut } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const loading = ref(true);
const saving = ref(false);
const editing = ref(false);
const error = ref('');
const success = ref('');
const fieldError = ref('');
const profile = ref(null);
const form = reactive({ name: '', position: '' });

const loadProfile = async () => {
  loading.value = true;
  error.value = '';
  try {
    profile.value = await apiGet('/api/me/profile');
  } catch (requestError) {
    error.value = requestError.message || 'Unable to load your profile.';
  } finally {
    loading.value = false;
  }
};

const beginEdit = () => {
  form.name = profile.value?.name || '';
  form.position = profile.value?.position || '';
  fieldError.value = '';
  success.value = '';
  editing.value = true;
};

const saveProfile = async () => {
  fieldError.value = '';
  error.value = '';
  success.value = '';
  if (!form.name.trim()) {
    fieldError.value = 'Name is required.';
    return;
  }
  saving.value = true;
  try {
    profile.value = await apiPut('/api/me/profile', { name: form.name, position: form.position });
    const user = JSON.parse(localStorage.getItem('authUser') || '{}');
    const updatedUser = { ...user, full_name: profile.value.name, name: profile.value.name, position: profile.value.position };
    localStorage.setItem('authUser', JSON.stringify(updatedUser));
    window.dispatchEvent(new CustomEvent('auth-user-updated', { detail: updatedUser }));
    editing.value = false;
    success.value = 'Profile updated successfully.';
  } catch (requestError) {
    error.value = requestError.message || 'Unable to save your profile.';
  } finally {
    saving.value = false;
  }
};

onMounted(loadProfile);
</script>

<style scoped>
.self-service-page { width: 100%; min-width: 0; }
.page-heading { display: flex; justify-content: space-between; align-items: center; gap: 14px; margin-bottom: 20px; }
.page-heading h2 { margin: 0 0 4px; color: #0f172a; font-size: 24px; }
.breadcrumb { margin: 0; color: #64748b; font-size: 13px; }
.profile-card { max-width: 760px; padding: 24px; border: 1px solid #e2e8f0; border-radius: 12px; background: white; }
.profile-details { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; margin: 0 0 24px; }
.profile-details div { min-width: 0; }
.profile-details dt { margin-bottom: 5px; color: #64748b; font-size: 12px; text-transform: uppercase; }
.profile-details dd { margin: 0; color: #18352b; font-weight: 600; overflow-wrap: anywhere; }
.profile-form { display: grid; gap: 14px; max-width: 480px; }
.profile-form label { display: grid; gap: 6px; color: #334155; font-size: 13px; font-weight: 600; }
.profile-form input { width: 100%; box-sizing: border-box; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 8px; }
.profile-form input:disabled { background: #f1f5f9; color: #64748b; }
.form-actions { display: flex; flex-wrap: wrap; gap: 10px; }
.primary-button, .secondary-button { display: inline-flex; align-items: center; justify-content: center; min-height: 38px; padding: 8px 14px; border: 0; border-radius: 8px; font-size: 13px; font-weight: 600; text-decoration: none; cursor: pointer; }
.primary-button { background: #047857; color: white; }
.secondary-button { border: 1px solid #cbd5e1; background: white; color: #334155; }
.api-error, .success-message, .field-error { font-size: 13px; }
.api-error { color: #b91c1c; }
.success-message { color: #047857; }
.field-error { margin: 0; color: #b91c1c; }
@media (max-width: 600px) { .page-heading { align-items: flex-start; flex-direction: column; } .profile-details { grid-template-columns: 1fr; } }
</style>
