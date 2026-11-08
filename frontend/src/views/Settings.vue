<template>
  <div class="settings-container">
    <div class="page-header">
      <div class="header-title">
        <h1>⚙️ System Settings & Configurations</h1>
        <p class="breadcrumb">Dashboard > Settings</p>
      </div>
    </div>

    <div class="settings-content">
      <div class="settings-card">
        <h3>Admin Profile Settings</h3>
        <EmptyState v-if="states.profile.loading" message="Loading profile..." />
        <template v-else>
          <p v-if="states.profile.error" role="alert" class="api-error">{{ states.profile.error }}</p>
          <p v-if="states.profile.success" role="status" class="success-message">{{ states.profile.success }}</p>
          <template v-if="profile">
            <div class="account-details">
              <p><strong>Name:</strong> {{ profile.name }}</p>
              <p><strong>Email:</strong> {{ profile.email }}</p>
              <p><strong>Employee ID:</strong> {{ profile.employee_id }}</p>
              <p><strong>Position:</strong> {{ profile.position || '—' }}</p>
              <p><strong>Department:</strong> {{ profile.department || '—' }}</p>
              <p><strong>Role:</strong> {{ profile.role }}</p>
            </div>
            <button v-if="!editingProfile" class="btn-save" type="button" @click="startProfileEdit">Edit</button>
            <form v-else class="settings-form" @submit.prevent="saveProfile">
              <label class="input-group">
                Name
                <input v-model="profileForm.name" class="form-input" type="text" />
                <span v-if="profileErrors.name" class="field-error">{{ profileErrors.name }}</span>
              </label>
              <label class="input-group">
                Email
                <input v-model="profileForm.email" class="form-input" type="email" />
                <span v-if="profileErrors.email" class="field-error">{{ profileErrors.email }}</span>
              </label>
              <label class="input-group">
                Position
                <input v-model="profileForm.position" class="form-input" type="text" />
              </label>
              <div class="form-actions">
                <button class="btn-save" type="submit" :disabled="states.profile.saving">{{ states.profile.saving ? 'Saving...' : 'Save Profile' }}</button>
                <button class="btn-cancel" type="button" @click="editingProfile = false">Cancel</button>
              </div>
            </form>

            <h4 class="subsection-title">Change Password</h4>
            <form class="settings-form" @submit.prevent="changePassword">
              <label class="input-group">
                Current Password
                <input v-model="passwordForm.current_password" class="form-input" type="password" autocomplete="current-password" />
                <span v-if="passwordErrors.current_password" class="field-error">{{ passwordErrors.current_password }}</span>
              </label>
              <label class="input-group">
                New Password
                <input v-model="passwordForm.new_password" class="form-input" type="password" autocomplete="new-password" />
                <span v-if="passwordErrors.new_password" class="field-error">{{ passwordErrors.new_password }}</span>
              </label>
              <label class="input-group">
                Confirm New Password
                <input v-model="passwordForm.confirm_password" class="form-input" type="password" autocomplete="new-password" />
                <span v-if="passwordErrors.confirm_password" class="field-error">{{ passwordErrors.confirm_password }}</span>
              </label>
              <button class="btn-save" type="submit" :disabled="states.profile.saving">{{ states.profile.saving ? 'Saving...' : 'Change Password' }}</button>
            </form>
          </template>
          <EmptyState v-else message="No signed-in account details available" />
        </template>
      </div>

      <template v-if="canReadSettings">
      <div class="settings-card">
        <h3>Company Profile Settings</h3>
        <EmptyState v-if="states.company.loading" message="Loading company settings..." />
        <template v-else>
          <p v-if="states.company.error" role="alert" class="api-error">{{ states.company.error }}</p>
          <p v-if="states.company.success" role="status" class="success-message">{{ states.company.success }}</p>
          <form class="settings-form" @submit.prevent="isAdmin && saveCompany()">
            <label class="input-group">
              Company Name
              <input v-model="company.company_name" class="form-input" type="text" :disabled="!isAdmin" />
              <span v-if="companyErrors.company_name" class="field-error">{{ companyErrors.company_name }}</span>
            </label>
            <label class="input-group">Address<input v-model="company.address" class="form-input" type="text" :disabled="!isAdmin" /></label>
            <label class="input-group">Phone<input v-model="company.phone" class="form-input" type="tel" :disabled="!isAdmin" /></label>
            <label class="input-group">
              Email
              <input v-model="company.email" class="form-input" type="email" :disabled="!isAdmin" />
              <span v-if="companyErrors.email" class="field-error">{{ companyErrors.email }}</span>
            </label>
            <label class="input-group">
              Website
              <input v-model="company.website" class="form-input" type="url" :disabled="!isAdmin" />
              <span v-if="companyErrors.website" class="field-error">{{ companyErrors.website }}</span>
            </label>
            <button v-if="isAdmin" class="btn-save" type="submit" :disabled="states.company.saving">{{ states.company.saving ? 'Saving...' : 'Save Company Settings' }}</button>
          </form>
        </template>
      </div>

      <div class="settings-card">
        <h3>Shift & Attendance Policy</h3>
        <EmptyState v-if="states.attendance.loading" message="Loading attendance policy..." />
        <template v-else>
          <p v-if="states.attendance.error" role="alert" class="api-error">{{ states.attendance.error }}</p>
          <p v-if="states.attendance.success" role="status" class="success-message">{{ states.attendance.success }}</p>
          <form class="settings-form" @submit.prevent="isAdmin && saveAttendance()">
            <label class="input-group">Work Start<input v-model="attendance.work_start" class="form-input" type="time" required :disabled="!isAdmin" /></label>
            <label class="input-group">Work End<input v-model="attendance.work_end" class="form-input" type="time" required :disabled="!isAdmin" /></label>
            <label class="input-group">Late After<input v-model="attendance.late_after" class="form-input" type="time" required :disabled="!isAdmin" /></label>
            <span v-if="attendanceErrors.times" class="field-error">{{ attendanceErrors.times }}</span>
            <fieldset class="day-fieldset">
              <legend>Working Days</legend>
              <label v-for="day in weekDays" :key="`work-${day}`" class="day-option">
                <input v-model="attendance.working_days" type="checkbox" :value="day" :disabled="!isAdmin" />
                {{ day }}
              </label>
            </fieldset>
            <fieldset class="day-fieldset">
              <legend>Weekly Off Days</legend>
              <label v-for="day in weekDays" :key="`off-${day}`" class="day-option">
                <input v-model="attendance.weekly_off_days" type="checkbox" :value="day" :disabled="!isAdmin" />
                {{ day }}
              </label>
            </fieldset>
            <button v-if="isAdmin" class="btn-save" type="submit" :disabled="states.attendance.saving">{{ states.attendance.saving ? 'Saving...' : 'Save Attendance Policy' }}</button>
          </form>
        </template>
      </div>

      <div class="settings-card">
        <h3>User Roles & Access Permissions</h3>
        <EmptyState v-if="states.roles.loading" message="Loading user roles..." />
        <template v-else>
          <p v-if="states.roles.error" role="alert" class="api-error">{{ states.roles.error }}</p>
          <p v-if="states.roles.success" role="status" class="success-message">{{ states.roles.success }}</p>
          <div class="roles-table-scroll">
            <table v-if="users.length" class="roles-table">
              <thead>
                <tr><th>Employee</th><th>Email</th><th>Role</th><th>Action</th></tr>
              </thead>
              <tbody>
                <tr v-for="user in users" :key="user.id">
                  <td><strong>{{ user.name }}</strong><br /><small>{{ user.employee_id }}</small></td>
                  <td>{{ user.email }}</td>
                  <td>
                    <select v-if="isAdmin" v-model="roleDrafts[user.id]" class="form-select" :aria-label="`Role for ${user.name}`">
                      <option v-for="role in roleDefinitions" :key="role.role" :value="role.role">{{ role.role }}</option>
                    </select>
                    <span v-else>{{ user.role }}</span>
                    <span v-if="roleErrors[user.id]" class="field-error">{{ roleErrors[user.id] }}</span>
                  </td>
                  <td><button v-if="isAdmin" class="btn-save row-save" type="button" :disabled="roleSavingId === user.id || roleDrafts[user.id] === user.role" @click="saveRole(user)">{{ roleSavingId === user.id ? 'Saving...' : 'Save' }}</button><span v-else>Read only</span></td>
                </tr>
              </tbody>
            </table>
            <EmptyState v-else message="No users available" />
          </div>
          <div class="permission-summary">
            <h4>Permissions by role</h4>
            <p v-for="role in roleDefinitions" :key="role.role"><strong>{{ role.role }}:</strong> {{ role.permissions.join(', ') }}</p>
            <EmptyState v-if="!roleDefinitions.length" message="No role permissions available" />
          </div>
        </template>
      </div>

      <div class="settings-card">
        <h3>System Preferences & Integrations</h3>
        <EmptyState v-if="states.preferences.loading" message="Loading preferences..." />
        <template v-else>
          <p v-if="states.preferences.error" role="alert" class="api-error">{{ states.preferences.error }}</p>
          <p v-if="states.preferences.success" role="status" class="success-message">{{ states.preferences.success }}</p>
          <form class="settings-form" @submit.prevent="isAdmin && savePreferences()">
            <label class="input-group">
              Currency
              <select v-model="preferences.currency" class="form-select" :disabled="!isAdmin">
                <option value="LKR">LKR</option>
                <option value="USD">USD</option>
                <option value="EUR">EUR</option>
                <option value="GBP">GBP</option>
              </select>
            </label>
            <label class="input-group">
              Date Format
              <select v-model="preferences.date_format" class="form-select" :disabled="!isAdmin">
                <option value="">Select date format</option>
                <option value="YYYY-MM-DD">YYYY-MM-DD</option>
                <option value="DD/MM/YYYY">DD/MM/YYYY</option>
                <option value="MM/DD/YYYY">MM/DD/YYYY</option>
              </select>
            </label>
            <label class="input-group">
              Timezone
              <select v-model="preferences.timezone" class="form-select" :disabled="!isAdmin">
                <option value="">Select timezone</option>
                <option v-for="timezone in timezones" :key="timezone" :value="timezone">{{ timezone }}</option>
              </select>
              <span v-if="preferenceErrors.timezone" class="field-error">{{ preferenceErrors.timezone }}</span>
            </label>
            <label class="toggle-item">
              <span>Email Notifications</span>
              <span class="switch">
                <input v-model="preferences.email_notifications" type="checkbox" :disabled="!isAdmin" />
                <span class="slider"></span>
              </span>
            </label>
            <button v-if="isAdmin" class="btn-save" type="submit" :disabled="states.preferences.saving">{{ states.preferences.saving ? 'Saving...' : 'Save Preferences' }}</button>
          </form>
        </template>
      </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { apiGet, apiPost, apiPut } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const currentUser = ref({});
const canReadSettings = computed(() => ['Admin', 'HR Manager'].includes(currentUser.value.role));
const isAdmin = computed(() => currentUser.value.role === 'Admin');
const profile = ref(null);
const editingProfile = ref(false);
const profileForm = reactive({ name: '', email: '', position: '' });
const passwordForm = reactive({ current_password: '', new_password: '', confirm_password: '' });
const profileErrors = ref({});
const passwordErrors = ref({});
const company = reactive({ company_name: '', address: '', phone: '', email: '', website: '' });
const companyErrors = ref({});
const attendance = reactive({ work_start: '', work_end: '', late_after: '', working_days: [], weekly_off_days: [] });
const attendanceErrors = ref({});
const preferences = reactive({ currency: 'LKR', date_format: '', timezone: 'Asia/Colombo', email_notifications: false });
const preferenceErrors = ref({});
const roleDefinitions = ref([]);
const users = ref([]);
const roleDrafts = reactive({});
const roleErrors = reactive({});
const roleSavingId = ref(null);
const weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
const timezones = ['Asia/Colombo', 'Asia/Kolkata', 'Asia/Singapore', 'Europe/London', 'UTC', 'America/New_York'];
const states = reactive({
  profile: { loading: true, saving: false, error: '', success: '' },
  company: { loading: true, saving: false, error: '', success: '' },
  attendance: { loading: true, saving: false, error: '', success: '' },
  roles: { loading: true, saving: false, error: '', success: '' },
  preferences: { loading: true, saving: false, error: '', success: '' },
});

const captureError = error => error instanceof Error ? error.message : 'Unexpected error. Please try again.';
const loadProfile = async () => {
  states.profile.loading = true;
  states.profile.error = '';
  try {
    profile.value = await apiGet('/api/profile');
    Object.assign(profileForm, { name: profile.value.name || '', email: profile.value.email || '', position: profile.value.position || '' });
    currentUser.value = { ...currentUser.value, ...profile.value, full_name: profile.value.name };
  } catch (error) {
    console.error('Unable to load profile', error);
    states.profile.error = captureError(error);
  } finally {
    states.profile.loading = false;
  }
};

const loadSettings = async () => {
  for (const section of ['company', 'attendance', 'preferences']) states[section].loading = true;
  try {
    const data = await apiGet('/api/settings');
    Object.assign(company, data.company || {});
    Object.assign(attendance, data.attendance || {});
    Object.assign(preferences, data.preferences || {});
  } catch (error) {
    console.error('Unable to load settings', error);
    for (const section of ['company', 'attendance', 'preferences']) states[section].error = captureError(error);
  } finally {
    for (const section of ['company', 'attendance', 'preferences']) states[section].loading = false;
  }
};

const loadRoles = async () => {
  states.roles.loading = true;
  states.roles.error = '';
  try {
    const [definitions, userList] = await Promise.all([apiGet('/api/roles'), apiGet('/api/users-roles')]);
    roleDefinitions.value = Array.isArray(definitions) ? definitions : [];
    users.value = Array.isArray(userList) ? userList : [];
    for (const user of users.value) roleDrafts[user.id] = user.role;
  } catch (error) {
    console.error('Unable to load roles', error);
    states.roles.error = captureError(error);
  } finally {
    states.roles.loading = false;
  }
};

const startProfileEdit = () => {
  Object.assign(profileForm, { name: profile.value.name || '', email: profile.value.email || '', position: profile.value.position || '' });
  profileErrors.value = {};
  states.profile.error = '';
  states.profile.success = '';
  editingProfile.value = true;
};

const saveProfile = async () => {
  profileErrors.value = {};
  states.profile.error = '';
  states.profile.success = '';
  if (!profileForm.name.trim()) profileErrors.value.name = 'Name is required.';
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profileForm.email.trim())) profileErrors.value.email = 'Enter a valid email address.';
  if (Object.keys(profileErrors.value).length) return;
  states.profile.saving = true;
  try {
    profile.value = await apiPut('/api/profile', {
      name: profileForm.name.trim(),
      email: profileForm.email.trim(),
      position: profileForm.position.trim(),
    });
    const savedUser = JSON.parse(localStorage.getItem('authUser') || '{}');
    const updatedUser = {
      ...savedUser,
      id: profile.value.id || savedUser.id,
      full_name: profile.value.name,
      name: profile.value.name,
      email: profile.value.email,
      position: profile.value.position,
      department: profile.value.department,
      role: profile.value.role,
      employee_id: profile.value.employee_id,
    };
    localStorage.setItem('authUser', JSON.stringify(updatedUser));
    currentUser.value = updatedUser;
    window.dispatchEvent(new CustomEvent('auth-user-updated', { detail: updatedUser }));
    editingProfile.value = false;
    states.profile.success = 'Profile updated successfully.';
  } catch (error) {
    console.error('Unable to save profile', error);
    states.profile.error = captureError(error);
  } finally {
    states.profile.saving = false;
  }
};

const changePassword = async () => {
  passwordErrors.value = {};
  states.profile.error = '';
  states.profile.success = '';
  if (!passwordForm.current_password) passwordErrors.value.current_password = 'Current password is required.';
  if (passwordForm.new_password.length < 8) passwordErrors.value.new_password = 'Use at least 8 characters.';
  if (passwordForm.new_password !== passwordForm.confirm_password) passwordErrors.value.confirm_password = 'Passwords do not match.';
  if (Object.keys(passwordErrors.value).length) return;
  states.profile.saving = true;
  try {
    const result = await apiPost('/api/profile/change-password', passwordForm);
    Object.assign(passwordForm, { current_password: '', new_password: '', confirm_password: '' });
    states.profile.success = result.message || 'Password changed successfully.';
  } catch (error) {
    console.error('Unable to change password', error);
    states.profile.error = captureError(error);
  } finally {
    states.profile.saving = false;
  }
};

const saveCompany = async () => {
  companyErrors.value = {};
  states.company.error = '';
  states.company.success = '';
  if (!company.company_name.trim()) companyErrors.value.company_name = 'Company name is required.';
  if (company.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(company.email.trim())) companyErrors.value.email = 'Enter a valid email address.';
  if (company.website && !/^https?:\/\/.+/i.test(company.website.trim())) companyErrors.value.website = 'Enter a valid website URL beginning with http:// or https://.';
  if (Object.keys(companyErrors.value).length) return;
  states.company.saving = true;
  try {
    Object.assign(company, await apiPut('/api/settings/company', company));
    states.company.success = 'Company settings saved successfully.';
  } catch (error) {
    console.error('Unable to save company settings', error);
    states.company.error = captureError(error);
  } finally {
    states.company.saving = false;
  }
};

const saveAttendance = async () => {
  attendanceErrors.value = {};
  states.attendance.error = '';
  states.attendance.success = '';
  const timePattern = /^(?:[01]\d|2[0-3]):[0-5]\d$/;
  if (![attendance.work_start, attendance.work_end, attendance.late_after].every(value => timePattern.test(value))) {
    attendanceErrors.value.times = 'Enter valid work start, work end, and late-after times.';
  } else if (attendance.work_end <= attendance.work_start) {
    attendanceErrors.value.times = 'Work end must be after work start.';
  }
  if (Object.keys(attendanceErrors.value).length) return;
  states.attendance.saving = true;
  try {
    Object.assign(attendance, await apiPut('/api/settings/attendance', attendance));
    states.attendance.success = 'Attendance policy saved successfully.';
  } catch (error) {
    console.error('Unable to save attendance settings', error);
    states.attendance.error = captureError(error);
  } finally {
    states.attendance.saving = false;
  }
};

const savePreferences = async () => {
  preferenceErrors.value = {};
  states.preferences.error = '';
  states.preferences.success = '';
  if (!preferences.timezone) preferenceErrors.value.timezone = 'Select a timezone.';
  if (Object.keys(preferenceErrors.value).length) return;
  states.preferences.saving = true;
  try {
    Object.assign(preferences, await apiPut('/api/settings/preferences', preferences));
    states.preferences.success = 'Preferences saved successfully.';
  } catch (error) {
    console.error('Unable to save preferences', error);
    states.preferences.error = captureError(error);
  } finally {
    states.preferences.saving = false;
  }
};

const saveRole = async user => {
  roleErrors[user.id] = '';
  states.roles.error = '';
  states.roles.success = '';
  roleSavingId.value = user.id;
  try {
    await apiPut(`/api/users-roles/${user.id}`, { role: roleDrafts[user.id] });
    await loadRoles();
    states.roles.success = `Role updated for ${user.name}.`;
  } catch (error) {
    console.error('Unable to update user role', error);
    roleErrors[user.id] = captureError(error);
  } finally {
    roleSavingId.value = null;
  }
};

onMounted(async () => {
  try {
    currentUser.value = JSON.parse(localStorage.getItem('authUser') || '{}');
  } catch (error) {
    console.error('Unable to read signed-in account', error);
    currentUser.value = {};
  }
  await loadProfile();
  if (canReadSettings.value) {
    await Promise.all([loadSettings(), loadRoles()]);
  } else {
    for (const section of ['company', 'attendance', 'preferences', 'roles']) states[section].loading = false;
  }
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

.settings-container {
  padding: 25px;
  font-family: 'Inter', sans-serif;
  background-color: #f8fafc;
  min-height: 100vh;
}
.page-header h1 { font-size: 24px; font-weight: 700; color: #062c22; margin: 0; }
.breadcrumb { font-size: 12px; color: #64748b; margin-top: 4px; }
.settings-content { display: flex; flex-direction: column; gap: 25px; margin-top: 25px; }

.settings-card {
  background: white; border-radius: 15px; padding: 25px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
.settings-card h3 {
  font-size: 16px;
  color: #062c22;
  margin: 0 0 20px;
  font-weight: 600;
  border-bottom: 2px solid #f1f5f9;
  padding-bottom: 10px;
}
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; }
.success-message { color: #047857; background: #ecfdf5; border: 1px solid #a7f3d0; border-radius: 8px; padding: 10px 14px; }
.account-details p { margin: 0 0 8px; color: #334155; font-size: 13px; }
.settings-form { display: flex; flex-direction: column; gap: 14px; max-width: 520px; }
.input-group { display: flex; flex-direction: column; gap: 6px; font-size: 12px; font-weight: 600; color: #475569; }
.form-input, .form-select {
  width: 100%;
  min-height: 40px;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  outline: none;
  background: #f8fafc;
  color: #334155;
  font-family: 'Inter', sans-serif;
}
.field-error { color: #b91c1c; font-size: 12px; font-weight: 400; }
.subsection-title { margin: 25px 0 14px; color: #062c22; font-size: 14px; }
.form-actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.btn-save {
  background: #10b981; color: white; border: none; padding: 10px 20px;
  border-radius: 25px; font-weight: 500; cursor: pointer; align-self: flex-start; margin-top: 4px;
  font-size: 13px;
}
.btn-save:hover { background: #059669; }
.btn-save:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-cancel { padding: 9px 18px; border: 1px solid #cbd5e1; border-radius: 25px; color: #475569; background: white; cursor: pointer; }
.day-fieldset { display: flex; flex-wrap: wrap; gap: 12px; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px; }
.day-fieldset legend { padding: 0 6px; color: #475569; font-size: 12px; font-weight: 600; }
.day-option { display: inline-flex; align-items: center; gap: 5px; color: #334155; font-size: 12px; }
.roles-table-scroll { width: 100%; overflow-x: auto; }
.roles-table { width: 100%; min-width: 600px; border-collapse: collapse; font-size: 13px; text-align: left; margin-top: 15px; }
.roles-table th { padding: 12px; background: #f8fafc; color: #475569; font-weight: 600; border-bottom: 1px solid #e2e8f0; }
.roles-table td { padding: 14px; border-bottom: 1px solid #f1f5f9; color: #334155; }
.roles-table small { color: #64748b; }
.row-save { padding: 7px 14px; }
.permission-summary { margin-top: 22px; color: #475569; font-size: 12px; }
.permission-summary h4 { margin: 0 0 10px; color: #062c22; font-size: 13px; }
.permission-summary p { margin: 7px 0; }
.toggle-item { display: flex; justify-content: space-between; align-items: center; color: #334155; font-size: 13px; font-weight: 600; }

.switch { position: relative; display: inline-block; width: 44px; height: 24px; }
.switch input { opacity: 0; width: 0; height: 0; }
.slider {
  position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0;
  background-color: #cbd5e1; transition: .3s; border-radius: 24px;
}
.slider:before {
  position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px;
  background-color: white; transition: .3s; border-radius: 50%;
}
input:checked + .slider { background-color: #10b981; }
input:checked + .slider:before { transform: translateX(20px); }
@media (max-width: 640px) {
  .settings-container { padding: 16px; }
  .settings-card { padding: 18px; }
}
</style>
