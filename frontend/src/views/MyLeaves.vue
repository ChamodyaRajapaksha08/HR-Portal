<template>
  <section class="self-service-page">
    <header class="page-heading">
      <div><h2>My Leaves</h2><p class="breadcrumb">Dashboard > My Leaves</p></div>
      <button class="primary-button apply-button" type="button" @click="openModal">+ Apply Leave</button>
    </header>
    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <p v-if="success" class="success-message" role="status">{{ success }}</p>
    <EmptyState v-if="loading" message="Loading your leave information..." />
    <template v-else>
      <div class="balance-grid">
        <article v-for="balance in balances" :key="balance.leave_type" class="balance-card">
          <span>{{ balance.leave_type }} Leave</span>
          <strong>{{ balance.remaining }} <small>days remaining</small></strong>
          <p>{{ balance.used }} used of {{ balance.entitlement }}</p>
        </article>
      </div>
      <section class="table-card">
        <h3>Leave Requests</h3>
        <div class="table-scroll">
          <table>
            <thead><tr><th>Leave Type</th><th>Date Range</th><th>Days</th><th>Reason</th><th>Status</th><th>Applied On</th></tr></thead>
            <tbody>
              <tr v-for="leave in leaves" :key="leave.id">
                <td>{{ leave.leave_type }}</td><td>{{ leave.from_date }} – {{ leave.to_date }}</td>
                <td>{{ leave.days }}</td><td>{{ leave.reason }}</td>
                <td><span class="status-badge" :class="statusClass(leave.status)">{{ leave.status }}</span></td>
                <td>{{ leave.applied_on }}</td>
              </tr>
              <tr v-if="!leaves.length"><td colspan="6"><EmptyState message="No leave requests yet" /></td></tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>

    <div v-if="modalOpen" class="modal-backdrop" @click.self="closeModal">
      <form class="leave-modal" @submit.prevent="submitLeave">
        <div class="modal-heading"><h3>Apply for Leave</h3><button type="button" aria-label="Close" @click="closeModal">×</button></div>
        <p v-if="formError" class="api-error" role="alert">{{ formError }}</p>
        <label>Leave Type
          <select v-model="form.leave_type">
            <option value="">Select leave type</option><option>Annual</option><option>Casual</option><option>Sick</option>
          </select>
          <span v-if="fieldErrors.leave_type" class="field-error">{{ fieldErrors.leave_type }}</span>
        </label>
        <div class="date-fields">
          <label>From<input v-model="form.from_date" type="date" /><span v-if="fieldErrors.from_date" class="field-error">{{ fieldErrors.from_date }}</span></label>
          <label>To<input v-model="form.to_date" type="date" /><span v-if="fieldErrors.to_date" class="field-error">{{ fieldErrors.to_date }}</span></label>
        </div>
        <div class="calculated-days">Requested days <strong>{{ calculatedDays }}</strong></div>
        <label>Reason<textarea v-model.trim="form.reason" rows="3"></textarea><span v-if="fieldErrors.reason" class="field-error">{{ fieldErrors.reason }}</span></label>
        <div class="modal-actions">
          <button class="primary-button" type="submit" :disabled="saving">{{ saving ? 'Submitting...' : 'Submit Request' }}</button>
          <button class="secondary-button" type="button" @click="closeModal">Cancel</button>
        </div>
      </form>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { apiGet, apiPost } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const loading = ref(true);
const saving = ref(false);
const error = ref('');
const success = ref('');
const formError = ref('');
const modalOpen = ref(false);
const balances = ref([]);
const leaves = ref([]);
const fieldErrors = ref({});
const form = reactive({ leave_type: '', from_date: '', to_date: '', reason: '' });
const calculatedDays = computed(() => {
  if (!form.from_date || !form.to_date) return 0;
  const from = Date.parse(`${form.from_date}T00:00:00Z`);
  const to = Date.parse(`${form.to_date}T00:00:00Z`);
  return Number.isNaN(from) || Number.isNaN(to) || to < from ? 0 : Math.floor((to - from) / 86400000) + 1;
});
const statusClass = status => (status || '').toLowerCase();

const loadData = async () => {
  loading.value = true;
  error.value = '';
  try {
    const [balanceData, leaveData] = await Promise.all([
      apiGet('/api/me/leave-balance'),
      apiGet('/api/me/leaves'),
    ]);
    balances.value = Array.isArray(balanceData.balances) ? balanceData.balances : [];
    leaves.value = Array.isArray(leaveData.requests) ? leaveData.requests : [];
  } catch (requestError) {
    error.value = requestError.message || 'Unable to load your leave information.';
  } finally {
    loading.value = false;
  }
};

const clearForm = () => {
  Object.assign(form, { leave_type: '', from_date: '', to_date: '', reason: '' });
  fieldErrors.value = {};
  formError.value = '';
};
const openModal = () => { clearForm(); modalOpen.value = true; };
const closeModal = () => { modalOpen.value = false; clearForm(); };

const submitLeave = async () => {
  fieldErrors.value = {};
  formError.value = '';
  success.value = '';
  if (!form.leave_type) fieldErrors.value.leave_type = 'Choose a leave type.';
  if (!form.from_date) fieldErrors.value.from_date = 'Choose a start date.';
  if (!form.to_date) fieldErrors.value.to_date = 'Choose an end date.';
  else if (form.from_date && form.to_date < form.from_date) fieldErrors.value.to_date = 'End date cannot be before start date.';
  if (!form.reason.trim()) fieldErrors.value.reason = 'Enter a reason for your request.';
  if (Object.keys(fieldErrors.value).length) return;

  saving.value = true;
  try {
    await apiPost('/api/me/leaves', {
      leave_type: form.leave_type,
      from_date: form.from_date,
      to_date: form.to_date,
      reason: form.reason.trim(),
    });
    closeModal();
    success.value = 'Leave request submitted successfully.';
    await loadData();
  } catch (requestError) {
    formError.value = requestError.message || 'Unable to submit your leave request.';
  } finally {
    saving.value = false;
  }
};

onMounted(loadData);
</script>

<style scoped>
.self-service-page { width: 100%; min-width: 0; }
.page-heading { display: flex; justify-content: space-between; align-items: center; gap: 14px; margin-bottom: 18px; }
.page-heading h2 { margin: 0 0 4px; color: #0f172a; font-size: 24px; }
.breadcrumb { margin: 0; color: #64748b; font-size: 13px; }
.primary-button, .secondary-button { display: inline-flex; align-items: center; justify-content: center; min-height: 38px; padding: 8px 14px; border: 0; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.primary-button { background: #047857; color: white; }
.secondary-button { border: 1px solid #cbd5e1; background: white; color: #334155; }
.balance-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 16px; margin-bottom: 20px; }
.balance-card, .table-card { padding: 18px; border: 1px solid #e2e8f0; border-radius: 12px; background: white; }
.balance-card > span { display: block; color: #64748b; font-size: 12px; font-weight: 600; text-transform: uppercase; }
.balance-card strong { display: block; margin-top: 8px; color: #18352b; font-size: 25px; }
.balance-card strong small { font-size: 12px; font-weight: 500; }
.balance-card p { margin: 6px 0 0; color: #64748b; font-size: 13px; }
.table-card h3 { margin: 0 0 14px; color: #18352b; font-size: 16px; }
.table-scroll { width: 100%; max-width: 100%; overflow-x: auto; }
table { width: 100%; border-collapse: collapse; text-align: left; }
th, td { padding: 12px 14px; border-bottom: 1px solid #e2e8f0; white-space: nowrap; font-size: 13px; }
th { color: #64748b; font-size: 11px; text-transform: uppercase; }
.status-badge { padding: 4px 8px; border-radius: 6px; background: #f1f5f9; color: #475569; font-size: 11px; font-weight: 600; }
.status-badge.approved { background: #dcfce7; color: #15803d; }
.status-badge.pending { background: #fef3c7; color: #b45309; }
.status-badge.rejected { background: #fee2e2; color: #b91c1c; }
.api-error, .field-error { color: #b91c1c; font-size: 13px; }
.success-message { color: #047857; font-size: 13px; }
.modal-backdrop { position: fixed; z-index: 1000; inset: 0; display: grid; place-items: center; padding: 18px; background: rgba(15, 23, 42, .5); }
.leave-modal { display: grid; gap: 14px; width: min(100%, 520px); max-height: calc(100vh - 36px); overflow-y: auto; padding: 22px; border-radius: 14px; background: white; box-shadow: 0 20px 50px rgba(15, 23, 42, .2); }
.modal-heading { display: flex; align-items: center; justify-content: space-between; }
.modal-heading h3 { margin: 0; color: #18352b; }
.modal-heading button { border: 0; background: none; color: #64748b; font-size: 26px; cursor: pointer; }
.leave-modal label { display: grid; gap: 6px; color: #334155; font-size: 13px; font-weight: 600; }
.leave-modal input, .leave-modal select, .leave-modal textarea { width: 100%; min-width: 0; box-sizing: border-box; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 8px; font: inherit; }
.date-fields { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.calculated-days { display: flex; justify-content: space-between; padding: 10px 12px; border-radius: 8px; background: #f0fdf4; color: #334155; font-size: 13px; }
.modal-actions { display: flex; flex-wrap: wrap; gap: 10px; }
@media (max-width: 600px) { .page-heading { align-items: flex-start; flex-wrap: wrap; } .apply-button { width: 100%; } .date-fields { grid-template-columns: 1fr; } }
</style>
