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
        <div class="form-grid">
          <div class="input-group">
            <label>System Admin Name</label>
            <input v-model="profile.name" type="text" class="form-input" />
          </div>
          <div class="input-group">
            <label>Admin Email Address</label>
            <input v-model="profile.email" type="email" class="form-input" />
          </div>
          <div class="input-group">
            <label>New Password</label>
            <input v-model="newPassword" type="password" placeholder="Password changes require backend support" class="form-input" />
          </div>
          <button class="btn-save" @click="saveProfile">Save Profile</button>
          <span v-if="messages.profile" role="status">{{ messages.profile }}</span>
        </div>
      </div>

      <div class="settings-card">
        <h3>Company Profile Settings</h3>
        <div class="form-grid dynamic-grid">
          <div class="input-group">
            <label>Organization Name</label>
            <input v-model="company.name" type="text" class="form-input" />
          </div>
          <div class="input-group">
            <label>Company Website</label>
            <input v-model="company.website" type="url" class="form-input" />
          </div>
          <div class="input-group">
            <label>Currency Format</label>
            <select v-model="company.currency" class="form-select">
              <option>Sri Lankan Rupee (LKR)</option>
              <option>US Dollar (USD)</option>
            </select>
          </div>
          <button class="btn-save" @click="saveCompany">Update Company info</button>
          <span v-if="messages.company" role="status">{{ messages.company }}</span>
        </div>
      </div>

      <div class="settings-card">
        <h3>Shift & Attendance Policy</h3>
        <div class="form-grid dynamic-grid">
          <div class="input-group">
            <label>Office Start Time (Late Threshold)</label>
            <input v-model="policy.startTime" type="time" class="form-input" />
            <span class="input-tip">If you arrive after this time, it will be recorded as 'Late'</span>
          </div>
          <div class="input-group">
            <label>Office End Time</label>
            <input v-model="policy.endTime" type="time" class="form-input" />
          </div>
          <div class="input-group">
            <label>Allowed Casual Leaves (Per Month)</label>
            <input v-model.number="policy.casualLeaves" type="number" min="0" class="form-input" />
          </div>
          <button class="btn-save" @click="savePolicy">Apply Policy</button>
          <span v-if="messages.policy" role="status">{{ messages.policy }}</span>
        </div>
      </div>

      <div class="settings-card">
        <h3>User Roles & Access Permissions</h3>
        <p class="section-desc">
Manage the permissions of other users using the system here</p>
        <table class="roles-table">
          <thead>
            <tr>
              <th>Role</th>
              <th>Dashboard</th>
              <th>Employees & Payroll</th>
              <th>System Settings</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><strong>Super Admin</strong></td>
              <td>🟢 Full Access</td>
              <td>🟢 Full Access</td>
              <td>🟢 Full Access</td>
              <td><span class="badge-active">Default</span></td>
            </tr>
            <tr>
              <td><strong>HR Manager</strong></td>
              <td>🟢 View Only</td>
              <td>🟢 Full Access</td>
              <td>❌ No Access</td>
              <td><span class="badge-active">Active</span></td>
            </tr>
            <tr>
              <td><strong>Employee (Staff)</strong></td>
              <td>❌ No Access</td>
              <td>ℹ️ View Own Payslip</td>
              <td>❌ No Access</td>
              <td><span class="badge-active">Active</span></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="settings-card">
        <h3>System Preferences & Integrations</h3>
        <div class="toggle-list">
          <div class="toggle-item">
            <div>
              <h4>Email Notifications</h4>
              <p>Send daily attendance logs automatically via email.</p>
            </div>
            <label class="switch">
              <input v-model="preferences.emailNotifications" type="checkbox" @change="savePreferences">
              <span class="slider"></span>
            </label>
          </div>

          <div class="toggle-item">
            <div>
              <h4>Slack Integration</h4>
              <p>Sync leave requests with company Slack channel.</p>
            </div>
            <label class="switch">
              <input v-model="preferences.slackIntegration" type="checkbox" @change="savePreferences">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';

const STORAGE_KEY = 'hrmSettings';
const profile = reactive({ name: 'Hashini Nisansala', email: 'admin@itum.com' });
const company = reactive({ name: 'ITUM HR Solutions', website: 'https://itum.mrt.ac.lk', currency: 'Sri Lankan Rupee (LKR)' });
const policy = reactive({ startTime: '08:30', endTime: '17:00', casualLeaves: 2 });
const preferences = reactive({ emailNotifications: true, slackIntegration: false });
const newPassword = ref('');
const messages = reactive({ profile: '', company: '', policy: '' });

onMounted(() => {
  const saved = localStorage.getItem(STORAGE_KEY);
  if (!saved) return;
  try {
    const settings = JSON.parse(saved);
    Object.assign(profile, settings.profile || {});
    Object.assign(company, settings.company || {});
    Object.assign(policy, settings.policy || {});
    Object.assign(preferences, settings.preferences || {});
  } catch (error) {
    console.error('Unable to read saved settings', error);
    window.alert('Saved settings could not be read. You can still update and save them.');
  }
});

const persistSettings = () => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({ profile, company, policy, preferences }));
};

const saveProfile = () => {
  try {
    persistSettings();
    messages.profile = newPassword.value
      ? 'Profile saved locally. Password changes need a backend endpoint and were not saved.'
      : 'Profile saved in this browser.';
    newPassword.value = '';
  } catch (error) {
    console.error('Unable to save profile settings', error);
    window.alert(`Unable to save profile settings: ${error.message}`);
  }
};

const saveCompany = () => {
  try {
    persistSettings();
    messages.company = 'Company settings saved in this browser.';
  } catch (error) {
    console.error('Unable to save company settings', error);
    window.alert(`Unable to save company settings: ${error.message}`);
  }
};

const savePolicy = () => {
  try {
    persistSettings();
    messages.policy = 'Attendance policy saved in this browser.';
  } catch (error) {
    console.error('Unable to save attendance policy', error);
    window.alert(`Unable to save attendance policy: ${error.message}`);
  }
};

const savePreferences = () => {
  try {
    persistSettings();
  } catch (error) {
    console.error('Unable to save preferences', error);
    window.alert(`Unable to save preferences: ${error.message}`);
  }
};
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

.section-desc { font-size: 13px; color: #64748b; margin-bottom: 15px; }

/* Form Layouts */
.form-grid { display: flex; flex-direction: column; gap: 15px; max-width: 450px; }
.input-group { display: flex; flex-direction: column; gap: 6px; }
.input-group label { font-size: 12px; font-weight: 600; color: #475569; }
.form-input, .form-select {
  padding: 10px 14px; border: 1px solid #e2e8f0; border-radius: 25px;
  font-size: 13px; outline: none; background: #f8fafc; color: #334155;
  font-family: 'Inter', sans-serif;
}
.input-tip { font-size: 11px; color: #047857; font-weight: 500; }

.btn-save {
  background: #10b981; color: white; border: none; padding: 10px 20px;
  border-radius: 25px; font-weight: 500; cursor: pointer; align-self: flex-start; margin-top: 10px;
  font-size: 13px;
}
.btn-save:hover { background: #059669; }

/* Roles Table Design */
.roles-table { width: 100%; border-collapse: collapse; font-size: 13px; text-align: left; margin-top: 15px; }
.roles-table th { padding: 12px; background: #f8fafc; color: #475569; font-weight: 600; border-bottom: 1px solid #e2e8f0; }
.roles-table td { padding: 14px; border-bottom: 1px solid #f1f5f9; color: #334155; }
.badge-active { background: #e0f2fe; color: #0369a1; padding: 3px 8px; border-radius: 12px; font-size: 11px; font-weight: 600; }

/* Toggle Switch Styles */
.toggle-list { display: flex; flex-direction: column; gap: 20px; }
.toggle-item { display: flex; justify-content: space-between; align-items: center; }
.toggle-item h4 { font-size: 14px; color: #062c22; margin: 0; font-weight: 600; }
.toggle-item p { font-size: 12px; color: #64748b; margin: 4px 0 0; }

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
</style>