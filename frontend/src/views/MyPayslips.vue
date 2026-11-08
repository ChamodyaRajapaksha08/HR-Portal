<template>
  <section class="payslips-page">
    <header class="page-heading"><div><h2>My Payslips</h2><p class="breadcrumb">Dashboard > My Payslips</p></div></header>
    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <EmptyState v-if="loading" message="Loading your payslips..." />
    <section v-else class="payslip-list">
      <button v-for="slip in payslips" :key="slip.id" class="payslip-row" type="button" @click="selected = slip">
        <span><strong>{{ payslipMonth(slip) }}</strong><small>Payroll #{{ slip.id }}</small></span>
        <span class="payslip-total">{{ money(slip.net_salary) }}</span>
        <span class="status-badge" :class="(slip.status || '').toLowerCase()">{{ slip.status }}</span>
        <span class="open-label">View payslip →</span>
      </button>
      <EmptyState v-if="!payslips.length" message="No payslips yet" />
    </section>

    <div v-if="selected" class="modal-backdrop" @click.self="selected = null">
      <section class="payslip-sheet">
        <div class="payslip-actions">
          <button type="button" class="secondary-button close-button" @click="selected = null">Close</button>
          <button type="button" class="primary-button" @click="printPayslip">Print</button>
        </div>
        <div class="print-content">
          <header class="receipt-header">
            <h1>PAYSLIP</h1>
            <p>{{ payslipMonth(selected) }}</p>
          </header>
          <div class="employee-meta">
            <p><span>Employee</span><strong>{{ user.name || user.full_name || user.email || '—' }}</strong></p>
            <p><span>Employee ID</span><strong>{{ user.employee_id || '—' }}</strong></p>
            <p><span>Payment Date</span><strong>{{ selected.paid_date || 'Pending' }}</strong></p>
          </div>
          <div class="amount-lines">
            <p><span>Basic Salary</span><strong>{{ money(selected.basic_salary) }}</strong></p>
            <p><span>Allowances</span><strong>{{ money(selected.allowances) }}</strong></p>
            <p><span>Deductions</span><strong>− {{ money(selected.deductions) }}</strong></p>
            <p class="net-line"><span>Net Salary</span><strong>{{ money(selected.net_salary) }}</strong></p>
          </div>
          <p class="receipt-status">Status: {{ selected.status }}</p>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { apiGet } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const loading = ref(true);
const error = ref('');
const payslips = ref([]);
const selected = ref(null);
const user = ref({});
const money = value => `LKR ${(Number(value) || 0).toLocaleString()}`;
const payslipMonth = slip => slip?.paid_date
  ? new Date(`${slip.paid_date}T00:00:00`).toLocaleDateString(undefined, { month: 'long', year: 'numeric' })
  : 'Unpaid payslip';
const printPayslip = () => window.print();

onMounted(async () => {
  try {
    user.value = JSON.parse(localStorage.getItem('authUser') || '{}');
    const data = await apiGet('/api/me/payslips');
    payslips.value = Array.isArray(data.payslips) ? data.payslips : [];
  } catch (requestError) {
    error.value = requestError.message || 'Unable to load your payslips.';
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.payslips-page { width: 100%; min-width: 0; }
.page-heading { margin-bottom: 20px; }
.page-heading h2 { margin: 0 0 4px; color: #0f172a; font-size: 24px; }
.breadcrumb { margin: 0; color: #64748b; font-size: 13px; }
.api-error { color: #b91c1c; }
.payslip-list { display: grid; gap: 12px; }
.payslip-row { display: flex; align-items: center; gap: 18px; width: 100%; padding: 16px 18px; border: 1px solid #e2e8f0; border-radius: 12px; background: white; text-align: left; cursor: pointer; }
.payslip-row > span:first-child { display: grid; gap: 4px; flex: 1; }
.payslip-row strong { color: #18352b; font-size: 14px; }
.payslip-row small { color: #64748b; font-size: 12px; }
.payslip-total { color: #18352b; font-weight: 700; }
.status-badge { padding: 4px 8px; border-radius: 6px; background: #f1f5f9; color: #475569; font-size: 11px; font-weight: 600; }
.status-badge.paid { background: #dcfce7; color: #15803d; }
.status-badge.pending { background: #fef3c7; color: #b45309; }
.open-label { color: #047857; font-size: 12px; font-weight: 600; }
.modal-backdrop { position: fixed; z-index: 1000; inset: 0; display: grid; place-items: center; padding: 18px; background: rgba(15, 23, 42, .55); }
.payslip-sheet { width: min(100%, 620px); max-height: calc(100vh - 36px); overflow: auto; padding: 24px; border-radius: 14px; background: white; box-shadow: 0 20px 50px rgba(15, 23, 42, .2); }
.payslip-actions { display: flex; justify-content: flex-end; gap: 10px; margin-bottom: 24px; }
.primary-button, .secondary-button { min-height: 38px; padding: 8px 14px; border: 0; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.primary-button { background: #047857; color: white; }
.secondary-button { border: 1px solid #cbd5e1; background: white; color: #334155; }
.receipt-header { padding-bottom: 18px; border-bottom: 2px solid #18352b; text-align: center; }
.receipt-header h1 { margin: 0; color: #18352b; letter-spacing: 2px; }
.receipt-header p { margin: 6px 0 0; color: #64748b; }
.employee-meta { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12px; padding: 18px 0; border-bottom: 1px solid #e2e8f0; }
.employee-meta p, .amount-lines p { display: grid; gap: 5px; margin: 0; }
.employee-meta span, .amount-lines span, .receipt-status { color: #64748b; font-size: 12px; }
.employee-meta strong { color: #334155; font-size: 13px; overflow-wrap: anywhere; }
.amount-lines { display: grid; gap: 14px; padding: 20px 0; }
.amount-lines p { display: flex; justify-content: space-between; gap: 12px; color: #334155; font-size: 14px; }
.amount-lines strong { color: #18352b; }
.net-line { padding-top: 15px; border-top: 1px solid #cbd5e1; font-size: 16px !important; font-weight: 700; }
.net-line strong { font-size: 18px; }
.receipt-status { text-align: right; }
@media (max-width: 600px) { .payslip-row { align-items: flex-start; flex-wrap: wrap; } .payslip-row > span:first-child { flex-basis: calc(100% - 20px); } .employee-meta { grid-template-columns: 1fr; } }
@media print {
  :global(body *) { visibility: hidden !important; }
  .modal-backdrop, .modal-backdrop * { visibility: visible !important; }
  .modal-backdrop { position: absolute; inset: 0; display: block; padding: 0; background: white; }
  .payslip-sheet { width: 100%; max-height: none; overflow: visible; padding: 0; border-radius: 0; box-shadow: none; }
  .payslip-actions { display: none !important; }
  .print-content { color: black; }
}
</style>
